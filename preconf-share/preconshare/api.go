package preconshare

import (
	"context"
	"errors"
	"time"

	"github.com/cairoeth/preconfirmations/preconf-share/jsonrpcserver"
	"github.com/cairoeth/preconfirmations/preconf-share/metrics"
	"github.com/cairoeth/preconfirmations/preconf-share/spike"
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

	bundleCacheSize = 1000
)

type SimScheduler interface {
	ScheduleRequest(ctx context.Context, bundle *SendRequestArgs, highPriority bool) error
}

type BundleStorage interface {
	GetBundleByMatchingHash(ctx context.Context, hash common.Hash) (*SendRequestArgs, error)
	InsertPreconf(ctx context.Context, preconf *ConfirmRequestArgs) error
	GetPreconfByMatchingHash(ctx context.Context, hash common.Hash) (*int64, *hexutil.Bytes, *int64, error)
	UpdatePreconfTimeBySignature(ctx context.Context, time int64, hash common.Hash) error
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

func (m *API) SendRequest(ctx context.Context, request SendRequestArgs) (_ SendRequestResponse, err error) {
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

	matchingHash, hasUnmatchedHash, err := ValidateRequest(&request, currentBlock, m.signer)
	if err != nil {
		logger.Warn("failed to validate request", zap.Error(err))
		return SendRequestResponse{}, err
	}
	m.knownBundleCache.Add(matchingHash, request)

	signerAddress := jsonrpcserver.GetSigner(ctx)
	origin := jsonrpcserver.GetOrigin(ctx)
	if request.Metadata == nil {
		request.Metadata = &RequestMetadata{}
	}
	request.Metadata.Signer = signerAddress
	request.Metadata.ReceivedAt = hexutil.Uint64(uint64(time.Now().UnixMicro()))
	request.Metadata.OriginID = origin
	request.Metadata.Prematched = !hasUnmatchedHash

	metrics.RecordBundleValidationDuration(time.Since(validateBundleTime).Milliseconds())

	metrics.IncSbundlesReceivedValid()
	highPriority := jsonrpcserver.GetPriority(ctx)
	err = m.scheduler.ScheduleRequest(ctx, &request, highPriority)
	if err != nil {
		metrics.IncRPCCallFailure(SendRequestEndpointName)
		logger.Error("Failed to schedule request simulation", zap.Error(err))
		return SendRequestResponse{}, ErrInternalServiceError
	}

	// sleep a bit to check if request was preconfirmed
	time.Sleep(1000 * time.Millisecond)

	// get preconf
	block, signature, _, err := m.bundleStorage.GetPreconfByMatchingHash(ctx, matchingHash)
	if err != nil {
		logger.Error("Failed to get preconfirmations", zap.Error(err))
		return SendRequestResponse{
			RequestHash: matchingHash,
			Signature:   nil,
			Block:       0,
		}, nil
	}

	_ = m.bundleStorage.UpdatePreconfTimeBySignature(ctx, time.Since(startAt).Milliseconds(), matchingHash)

	return SendRequestResponse{
		RequestHash: matchingHash,
		Signature:   signature,
		Block:       hexutil.Uint64(uint64(*block)),
	}, nil
}

func (m *API) ConfirmRequest(ctx context.Context, confirmation ConfirmRequestArgs) (_ ConfirmRequestResponse, err error) {
	logger := m.log
	startAt := time.Now()
	defer func() {
		metrics.RecordRPCCallDuration(ConfirmRequestEndpointName, time.Since(startAt).Milliseconds())
	}()

	logger.Info("Received request connfirmation")

	// TODO: validate confirmation (check signature, that is targeted for the current request, etc.)

	// Store preconfirmation in database.
	err = m.bundleStorage.InsertPreconf(ctx, &confirmation)
	if err != nil {
		metrics.IncRPCCallFailure(ConfirmRequestEndpointName)
		logger.Error("Failed to insert confirmation in db", zap.Error(err))
		return ConfirmRequestResponse{}, ErrInternalServiceError
	}

	return ConfirmRequestResponse{
		Valid: true,
	}, nil
}

func (m *API) GetRequest(ctx context.Context, request GetRequestArgs) (_ GetRequestResponse, err error) {
	logger := m.log
	startAt := time.Now()
	defer func() {
		metrics.RecordRPCCallDuration(GetRequestEndpointName, time.Since(startAt).Milliseconds())
	}()

	// Get preconfirmation in database.
	block, signature, time, err := m.bundleStorage.GetPreconfByMatchingHash(ctx, request.Hash)
	if err != nil {
		logger.Error("Failed to get preconfirmations", zap.Error(err))
		return GetRequestResponse{
			Signature: nil,
			Block:     0,
			Time:      0,
		}, nil
	}

	return GetRequestResponse{
		Signature: signature,
		Block:     hexutil.Uint64(uint64(*block)),
		Time:      hexutil.Uint64(uint64(*time)),
	}, nil
}
