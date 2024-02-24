package preconshare

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/cairoeth/preconfirmations-avs/simqueue"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"go.uber.org/zap"
)

// SimulationResult is responsible for processing simulation results
// NOTE: That error should be returned only if simulation should be retried, for example if redis is down
type SimulationResult interface {
	SimulatedBundle(ctx context.Context, args *SendMevBundleArgs, info simqueue.QueueItemInfo) error
}

type Storage interface {
	InsertBundleForStats(ctx context.Context, bundle *SendMevBundleArgs) (known bool, err error)
	InsertBundleForBuilder(ctx context.Context, bundle *SendMevBundleArgs, result *SimMevBundleResponse, targetBlock uint64) error
	InsertHistoricalHint(ctx context.Context, currentBlock uint64, hint *Hint) error
}

type SimulationResultBackend struct {
	log              *zap.Logger
	hint             HintBackend
	eth              EthClient
	store            Storage
}

func NewSimulationResultBackend(log *zap.Logger, hint HintBackend, eth EthClient, store Storage) *SimulationResultBackend {
	return &SimulationResultBackend{
		log:              log,
		hint:             hint,
		eth:              eth,
		store:            store,
	}
}

func izZeroPriorityFeeTX(bundle *SendMevBundleArgs) bool {
	if len(bundle.Body) != 1 {
		return false // not a single tx bundle
	}
	if bundle.Body[0].Bundle != nil {
		return false // bundle, not a single tx bundle
	}
	var tx types.Transaction
	btx := bundle.Body[0]
	if btx.Tx == nil {
		return false // incorrect bundle
	}

	err := tx.UnmarshalBinary(*btx.Tx)
	if err != nil {
		return false // incorrect bundle
	}

	return tx.GasTipCap().Cmp(big.NewInt(0)) == 0
}

// SimulatedBundle is called when simulation is done
// NOTE: we return error only if we want to retry the simulation
func (s *SimulationResultBackend) SimulatedBundle(ctx context.Context,
	bundle *SendMevBundleArgs, _ simqueue.QueueItemInfo,
) error {
	start := time.Now()

	var hash common.Hash
	if bundle.Metadata != nil {
		hash = bundle.Metadata.BundleHash
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

		// if sim.Success {
		if !knownBundle {
			start := time.Now()
			err := s.ProcessHints(ctx, bundle)
			logger.Debug("Processed hints", zap.Duration("duration", time.Since(start)), zap.Error(err))
			if err != nil {
				logger.Error("Failed to process hints", zap.Error(err))
			}
		}

		// 	start = time.Now()
		// 	err = s.store.InsertBundleForBuilder(ctx, bundle, sim, uint64(sim.StateBlock)+1)
		// 	logger.Debug("Inserted bundle for builder", zap.Duration("duration", time.Since(start)), zap.Error(err))
		// 	if err != nil {
		// 		logger.Error("Failed to insert bundle for builder", zap.Error(err))
		// 	}
		// }
	}()

	// check bundle mev priority fee
	isZeroFee := izZeroPriorityFeeTX(bundle)
	if isZeroFee {
		logger.Debug("Bundle has zero priority fee, skipping builder")
	}

	// failed bundles are not sent to builders, we also don't send single-tx bundles with zero priority fee
	// if sim.Success && !isZeroFee {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		s.builders.SendBundle(ctx, logger, bundle, uint64(sim.StateBlock)+1)
	// 	}()
	// }

	wg.Wait()
	log.Info("Bundle processed", zap.String("bundle", hash.Hex()), zap.Duration("duration", time.Since(start)))
	return nil
}

func (s *SimulationResultBackend) ProcessHints(ctx context.Context, bundle *SendMevBundleArgs) error {
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
