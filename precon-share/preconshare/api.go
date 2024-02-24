package preconshare

import (
	"context"
	"errors"
	"time"

	"github.com/cairoeth/preconfirmations-avs/precon-share/jsonrpcserver"
	"github.com/cairoeth/preconfirmations-avs/precon-share/metrics"
	"github.com/cairoeth/preconfirmations-avs/precon-share/spike"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
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
	ScheduleBundleSimulation(ctx context.Context, bundle *SendMevBundleArgs, highPriority bool) error
}

type BundleStorage interface {
	GetBundleByMatchingHash(ctx context.Context, hash common.Hash) (*SendMevBundleArgs, error)
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

type API struct {
	log *zap.Logger

	scheduler      SimScheduler
	bundleStorage  BundleStorage
	eth            EthClient
	signer         types.Signer
	simBackends    []SimulationBackend
	simRateLimiter *rate.Limiter

	spikeManager      *spike.Manager[*SendMevBundleArgs]
	knownBundleCache  *lru.Cache[common.Hash, SendMevBundleArgs]
	cancellationCache *RedisCancellationCache
}

func NewAPI(
	log *zap.Logger,
	scheduler SimScheduler, bundleStorage BundleStorage, eth EthClient, signer types.Signer,
	simBackends []SimulationBackend, simRateLimit rate.Limit, cancellationCache *RedisCancellationCache,
	sbundleValidDuration time.Duration,
) *API {
	sm := spike.NewManager(func(ctx context.Context, k string) (*SendMevBundleArgs, error) {
		return bundleStorage.GetBundleByMatchingHash(ctx, common.HexToHash(k))
	}, sbundleValidDuration)

	return &API{
		log: log,

		scheduler:         scheduler,
		bundleStorage:     bundleStorage,
		eth:               eth,
		signer:            signer,
		simBackends:       simBackends,
		simRateLimiter:    rate.NewLimiter(simRateLimit, 1),
		spikeManager:      sm,
		knownBundleCache:  lru.NewCache[common.Hash, SendMevBundleArgs](bundleCacheSize),
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

func (m *API) SendBundle(ctx context.Context, bundle SendMevBundleArgs) (_ SendMevBundleResponse, err error) {
	logger := m.log
	startAt := time.Now()
	defer func() {
		metrics.RecordRPCCallDuration(SendBundleEndpointName, time.Since(startAt).Milliseconds())
	}()
	metrics.IncSbundlesReceived()

	validateBundleTime := time.Now()
	currentBlock, err := m.eth.BlockNumber(ctx)
	if err != nil {
		metrics.IncRPCCallFailure(SendBundleEndpointName)
		logger.Error("failed to get current block", zap.Error(err))
		return SendMevBundleResponse{}, ErrInternalServiceError
	}

	hash, hasUnmatchedHash, err := ValidateBundle(&bundle, currentBlock, m.signer)
	if err != nil {
		logger.Warn("failed to validate bundle", zap.Error(err))
		return SendMevBundleResponse{}, err
	}
	if oldBundle, ok := m.knownBundleCache.Get(hash); ok {
		if !newerInclusion(&oldBundle, &bundle) {
			logger.Debug("bundle already known, ignoring", zap.String("hash", hash.Hex()))
			return SendMevBundleResponse{hash}, nil
		}
	}
	m.knownBundleCache.Add(hash, bundle)

	signerAddress := jsonrpcserver.GetSigner(ctx)
	origin := jsonrpcserver.GetOrigin(ctx)
	if bundle.Metadata == nil {
		bundle.Metadata = &MevBundleMetadata{}
	}
	bundle.Metadata.Signer = signerAddress
	bundle.Metadata.ReceivedAt = hexutil.Uint64(uint64(time.Now().UnixMicro()))
	bundle.Metadata.OriginID = origin
	bundle.Metadata.Prematched = !hasUnmatchedHash

	metrics.RecordBundleValidationDuration(time.Since(validateBundleTime).Milliseconds())

	if hasUnmatchedHash {
		var unmatchedHash common.Hash
		if len(bundle.Body) > 0 && bundle.Body[0].Hash != nil {
			unmatchedHash = *bundle.Body[0].Hash
		} else {
			return SendMevBundleResponse{}, ErrInternalServiceError
		}
		fetchUnmatchedTime := time.Now()
		unmatchedBundle, err := m.spikeManager.GetResult(ctx, unmatchedHash.String())
		metrics.RecordBundleFetchUnmatchedDuration(time.Since(fetchUnmatchedTime).Milliseconds())
		if err != nil {
			logger.Error("Failed to fetch unmatched bundle", zap.Error(err), zap.String("matching_hash", unmatchedHash.Hex()))
			metrics.IncRPCCallFailure(SendBundleEndpointName)
			return SendMevBundleResponse{}, ErrBackrunNotFound
		}
		if privacy := unmatchedBundle.Privacy; privacy == nil && privacy.Hints.HasHint(HintHash) {
			// if the unmatched bundle have not configured privacy or has not set the hash hint
			// then we cannot backrun it
			return SendMevBundleResponse{}, ErrBackrunInvalidBundle
		}
		bundle.Body[0].Bundle = unmatchedBundle
		bundle.Body[0].Hash = nil
		// replace matching hash with actual bundle hash
		findAndReplace(bundle.Metadata.BodyHashes, unmatchedHash, unmatchedBundle.Metadata.BundleHash)
		// send 90 % of the refund to the unmatched bundle or the suggested refund if set
		refundPercent := RefundPercent
		if unmatchedBundle.Privacy != nil && unmatchedBundle.Privacy.WantRefund != nil {
			refundPercent = *unmatchedBundle.Privacy.WantRefund
		}
		bundle.Validity.Refund = []RefundConstraint{{0, refundPercent}}
		MergePrivacyBuilders(&bundle)
		err = MergeInclusionIntervals(&bundle.Inclusion, &unmatchedBundle.Inclusion)
		if err != nil {
			return SendMevBundleResponse{}, ErrBackrunInclusion
		}
	}

	metrics.IncSbundlesReceivedValid()
	highPriority := jsonrpcserver.GetPriority(ctx)
	err = m.scheduler.ScheduleBundleSimulation(ctx, &bundle, highPriority)
	if err != nil {
		metrics.IncRPCCallFailure(SendBundleEndpointName)
		logger.Error("Failed to schedule bundle simulation", zap.Error(err))
		return SendMevBundleResponse{}, ErrInternalServiceError
	}

	return SendMevBundleResponse{
		BundleHash: hash,
	}, nil
}
