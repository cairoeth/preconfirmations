package preconshare

import (
	"context"
	"sync"
	"time"

	"github.com/cairoeth/preconfirmations/preconf-share/simqueue"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

// SimulationResult is responsible for processing simulation results
// NOTE: That error should be returned only if simulation should be retried, for example if redis is down
type SimulationResult interface {
	SimulatedBundle(ctx context.Context, args *SendRequestArgs, info simqueue.QueueItemInfo) error
}

type Storage interface {
	InsertBundleForStats(ctx context.Context, bundle *SendRequestArgs) (known bool, err error)
	InsertHistoricalHint(ctx context.Context, currentBlock uint64, hint *Hint) error
	GetPreconfByMatchingHash(ctx context.Context, hash common.Hash) (*int64, *hexutil.Bytes, *int64, error)
}

type SimulationResultBackend struct {
	log   *zap.Logger
	hint  HintBackend
	eth   EthClient
	store Storage
}

func NewSimulationResultBackend(log *zap.Logger, hint HintBackend, eth EthClient, store Storage) *SimulationResultBackend {
	return &SimulationResultBackend{
		log:   log,
		hint:  hint,
		eth:   eth,
		store: store,
	}
}

// SimulatedBundle is called when simulation is done
// NOTE: we return error only if we want to retry the simulation
func (s *SimulationResultBackend) SimulatedBundle(ctx context.Context,
	bundle *SendRequestArgs, _ simqueue.QueueItemInfo,
) error {
	start := time.Now()

	var hash common.Hash
	if bundle.Metadata != nil {
		hash = bundle.Metadata.MatchingHash
	}
	logger := s.log.With(zap.String("bundle", hash.Hex()))

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		start := time.Now()
		knownBundle, err := s.store.InsertBundleForStats(ctx, bundle)
		logger.Debug("Inserted bundle for stats", zap.Duration("duration", time.Since(start)), zap.Error(err))
		if err != nil {
			logger.Error("Failed to insert bundle for stats", zap.Error(err))
		}

		if !knownBundle {
			start := time.Now()
			err := s.ProcessHints(ctx, bundle)
			logger.Debug("Processed hints", zap.Duration("duration", time.Since(start)), zap.Error(err))
			if err != nil {
				logger.Error("Failed to process hints", zap.Error(err))
			}
		}

		logger.Debug("Bundle processed, waiting 1000ms for preconfs")

		// sleep 1000ms to receive preconfirmations
		time.Sleep(1000 * time.Millisecond)

		logger.Debug("Wait over, checking preconfs received")

		// Fetch preconfirmations from database that match the request hash
		block, signature, _, err := s.store.GetPreconfByMatchingHash(ctx, hash)
		if err != nil {
			logger.Error("Failed to get preconfirmations", zap.Error(err))
			return
		}

		logger.Info("Preconfirmation found", zap.String("signature", common.Bytes2Hex(*signature)), zap.Uint64("block", uint64(*block)))

		// Sending signed transactions to the operator
		rpcClient := jsonrpc.NewClient("http://localhost:8000/receive")
		_, err = rpcClient.Call(ctx, "receive", bundle.Body)
		if err != nil {
			logger.Error("Failed to send signed transactions", zap.Error(err))
			return
		}
	}()

	logger.Debug("Request finalized", zap.String("bundle", hash.Hex()), zap.Duration("duration", time.Since(start)))

	wg.Wait()
	return nil
}

func (s *SimulationResultBackend) ProcessHints(ctx context.Context, bundle *SendRequestArgs) error {
	if bundle.Privacy == nil {
		return nil
	}
	if !bundle.Privacy.Hints.HasHint(HintHash) {
		return nil
	}

	extractedHints, err := ExtractHints(bundle)
	if err != nil {
		return err
	}

	s.log.Info("Notifying hint", zap.String("hash", extractedHints.Hash.Hex()))
	err = s.hint.NotifyHint(ctx, &extractedHints)
	if err != nil {
		return err
	}

	block, err := s.eth.BlockNumber(ctx)
	if err != nil {
		return err
	}
	err = s.store.InsertHistoricalHint(ctx, block, &extractedHints)
	if err != nil {
		return err
	}

	return nil
}
