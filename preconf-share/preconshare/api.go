package preconshare

import (
	"context"
	"errors"
	"time"

	"github.com/cairoeth/preconfirmations-avs/preconf-share/jsonrpcserver"
	"github.com/cairoeth/preconfirmations-avs/preconf-share/metrics"
	"github.com/cairoeth/preconfirmations-avs/preconf-share/spike"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

var (
	ErrInvalidInclusion      = errors.New("invalid inclusion")
	ErrInvalidBundleBodySize = errors.New("invalid bundle body size")
	ErrInvalidBundleBody     = errors.New("invalid bundle body")
	ErrBackrunNotFound       = errors.New("backrun not found")
	ErrBackrunInvalidBundle  = errors.New("backrun invalid bundle")
	ErrBackrunInclusion      = errors.New("backrun invalid inclusion")

	ErrInternalServiceError = errors.New("mev-share service error")

	simBundleTimeout    = 500 * time.Millisecond
	cancelBundleTimeout = 3 * time.Second
	bundleCacheSize     = 1000
)

type SimScheduler interface {
	ScheduleBundleSimulation(ctx context.Context, bundle *SendRequestArgs, highPriority bool) error
}

type BundleStorage interface {
	GetBundleByMatchingHash(ctx context.Context, hash common.Hash) (*SendRequestArgs, error)
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type API struct {
	log *zap.Logger

	scheduler     SimScheduler
	bundleStorage BundleStorage
	eth           EthClient
	signer        types.Signer
	simBackends   []SimulationBackend

	spikeManager      *spike.Manager[*SendRequestArgs]
	knownBundleCache  *lru.Cache[common.Hash, SendRequestArgs]
	cancellationCache *RedisCancellationCache
}

func NewAPI(
	log *zap.Logger,
	scheduler SimScheduler, bundleStorage BundleStorage, eth EthClient, signer types.Signer,
	simBackends []SimulationBackend, cancellationCache *RedisCancellationCache,
	sbundleValidDuration time.Duration,
) *API {
	sm := spike.NewManager(func(ctx context.Context, k string) (*SendRequestArgs, error) {
		return bundleStorage.GetBundleByMatchingHash(ctx, common.HexToHash(k))
	}, sbundleValidDuration)

	return &API{
		log: log,

		scheduler:         scheduler,
		bundleStorage:     bundleStorage,
		eth:               eth,
		signer:            signer,
		simBackends:       simBackends,
		spikeManager:      sm,
		knownBundleCache:  lru.NewCache[common.Hash, SendRequestArgs](bundleCacheSize),
		cancellationCache: cancellationCache,
	}
}

func findAndReplace(strs []common.Hash, old, replacer common.Hash) bool {
	var found bool
	for i, str := range strs {
		if str == old {
			strs[i] = replacer
			found = true
		}
	}
	return found
}

func (m *API) SendRequest(ctx context.Context, bundle SendRequestArgs) (_ SendRequestResponse, err error) {
	logger := m.log
	startAt := time.Now()
	defer func() {
		metrics.RecordRPCCallDuration(SendRequestEndpointName, time.Since(startAt).Milliseconds())
	}()
	metrics.IncSbundlesReceived()

	validateBundleTime := time.Now()
	currentBlock, err := m.eth.BlockNumber(ctx)
	if err != nil {
		metrics.IncRPCCallFailure(SendRequestEndpointName)
		logger.Error("failed to get current block", zap.Error(err))
		return SendRequestResponse{}, ErrInternalServiceError
	}

	hash, hasUnmatchedHash, err := ValidateBundle(&bundle, currentBlock, m.signer)
	if err != nil {
		logger.Warn("failed to validate bundle", zap.Error(err))
		return SendRequestResponse{}, err
	}
	if oldBundle, ok := m.knownBundleCache.Get(hash); ok {
		if !newerInclusion(&oldBundle, &bundle) {
			logger.Debug("bundle already known, ignoring", zap.String("hash", hash.Hex()))
			return SendRequestResponse{hash}, nil
		}
	}
	m.knownBundleCache.Add(hash, bundle)

	signerAddress := jsonrpcserver.GetSigner(ctx)
	origin := jsonrpcserver.GetOrigin(ctx)
	if bundle.Metadata == nil {
		bundle.Metadata = &RequestMetadata{}
	}
	bundle.Metadata.Signer = signerAddress
	bundle.Metadata.ReceivedAt = hexutil.Uint64(uint64(time.Now().UnixMicro()))
	bundle.Metadata.OriginID = origin
	bundle.Metadata.Prematched = !hasUnmatchedHash

	metrics.RecordBundleValidationDuration(time.Since(validateBundleTime).Milliseconds())

	// if hasUnmatchedHash {
	// 	var unmatchedHash common.Hash
	// 	if len(bundle.Body) > 0 && bundle.Body[0].Hash != nil {
	// 		unmatchedHash = *bundle.Body[0].Hash
	// 	} else {
	// 		return SendRequestResponse{}, ErrInternalServiceError
	// 	}
	// 	fetchUnmatchedTime := time.Now()
	// 	unmatchedBundle, err := m.spikeManager.GetResult(ctx, unmatchedHash.String())
	// 	metrics.RecordBundleFetchUnmatchedDuration(time.Since(fetchUnmatchedTime).Milliseconds())
	// 	if err != nil {
	// 		logger.Error("Failed to fetch unmatched bundle", zap.Error(err), zap.String("matching_hash", unmatchedHash.Hex()))
	// 		metrics.IncRPCCallFailure(SendRequestEndpointName)
	// 		return SendRequestResponse{}, ErrBackrunNotFound
	// 	}
	// 	if privacy := unmatchedBundle.Privacy; privacy == nil && privacy.Hints.HasHint(HintHash) {
	// 		// if the unmatched bundle have not configured privacy or has not set the hash hint
	// 		// then we cannot backrun it
	// 		return SendRequestResponse{}, ErrBackrunInvalidBundle
	// 	}
	// 	bundle.Body[0].Bundle = unmatchedBundle
	// 	bundle.Body[0].Hash = nil
	// 	// replace matching hash with actual bundle hash
	// 	findAndReplace(bundle.Metadata.BodyHashes, unmatchedHash, unmatchedBundle.Metadata.BundleHash)
	// 	// send 90 % of the refund to the unmatched bundle or the suggested refund if set
	// 	refundPercent := RefundPercent
	// 	if unmatchedBundle.Privacy != nil && unmatchedBundle.Privacy.WantRefund != nil {
	// 		refundPercent = *unmatchedBundle.Privacy.WantRefund
	// 	}
	// 	bundle.Validity.Refund = []RefundConstraint{{0, refundPercent}}
	// 	MergePrivacyBuilders(&bundle)
	// 	err = MergeInclusionIntervals(&bundle.Inclusion, &unmatchedBundle.Inclusion)
	// 	if err != nil {
	// 		return SendRequestResponse{}, ErrBackrunInclusion
	// 	}
	// }

	metrics.IncSbundlesReceivedValid()
	highPriority := jsonrpcserver.GetPriority(ctx)
	err = m.scheduler.ScheduleBundleSimulation(ctx, &bundle, highPriority)
	if err != nil {
		metrics.IncRPCCallFailure(SendRequestEndpointName)
		logger.Error("Failed to schedule bundle simulation", zap.Error(err))
		return SendRequestResponse{}, ErrInternalServiceError
	}

	return SendRequestResponse{
		BundleHash: hash,
	}, nil
}

func (m *API) ConfirmRequest(ctx context.Context, bundle SendRequestArgs) {
	// logger := m.log
	// startAt := time.Now()
	// defer func() {
	// 	metrics.RecordRPCCallDuration(ConfirmRequestEndpointName, time.Since(startAt).Milliseconds())
	// }()
	// metrics.IncSbundlesReceived()

	// validateBundleTime := time.Now()
	// currentBlock, err := m.eth.BlockNumber(ctx)
	// if err != nil {
	// 	metrics.IncRPCCallFailure(ConfirmRequestEndpointName)
	// 	logger.Error("failed to get current block", zap.Error(err))
	// 	return SendRequestResponse{}, ErrInternalServiceError
	// }

	// hash, hasUnmatchedHash, err := ValidateBundle(&bundle, currentBlock, m.signer)
	// if err != nil {
	// 	logger.Warn("failed to validate bundle", zap.Error(err))
	// 	return SendRequestResponse{}, err
	// }
	// if oldBundle, ok := m.knownBundleCache.Get(hash); ok {
	// 	if !newerInclusion(&oldBundle, &bundle) {
	// 		logger.Debug("bundle already known, ignoring", zap.String("hash", hash.Hex()))
	// 		return SendRequestResponse{hash}, nil
	// 	}
	// }
	// m.knownBundleCache.Add(hash, bundle)

	// signerAddress := jsonrpcserver.GetSigner(ctx)
	// origin := jsonrpcserver.GetOrigin(ctx)
	// if bundle.Metadata == nil {
	// 	bundle.Metadata = &RequestMetadata{}
	// }
	// bundle.Metadata.Signer = signerAddress
	// bundle.Metadata.ReceivedAt = hexutil.Uint64(uint64(time.Now().UnixMicro()))
	// bundle.Metadata.OriginID = origin
	// bundle.Metadata.Prematched = !hasUnmatchedHash

	// metrics.RecordBundleValidationDuration(time.Since(validateBundleTime).Milliseconds())

	// metrics.IncSbundlesReceivedValid()
	// highPriority := jsonrpcserver.GetPriority(ctx)
	// err = m.scheduler.ScheduleBundleSimulation(ctx, &bundle, highPriority)
	// if err != nil {
	// 	metrics.IncRPCCallFailure(ConfirmRequestEndpointName)
	// 	logger.Error("Failed to schedule bundle simulation", zap.Error(err))
	// 	return SendRequestResponse{}, ErrInternalServiceError
	// }

	// return SendRequestResponse{
	// 	BundleHash: hash,
	// }, nil
}
