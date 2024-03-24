package preconshare

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/cairoeth/preconfirmations/preconf-share/metrics"
	"github.com/cairoeth/preconfirmations/preconf-share/simqueue"
	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

var (
	consumeSimulationTimeout = 15 * time.Second
	simCacheTimeout          = 1 * time.Second
)

type SimQueue struct {
	log            *zap.Logger
	queue          simqueue.Queue
	eth            EthClient
	workers        []SimulationWorker
	workersPerNode int
}

func NewQueue(
	log *zap.Logger, queue simqueue.Queue, eth EthClient, sim []SimulationBackend, simRes SimulationResult,
	workersPerNode int, backgroundWg *sync.WaitGroup, cancelCache *RedisCancellationCache,
) *SimQueue {
	log = log.Named("queue")
	q := &SimQueue{
		log:            log,
		queue:          queue,
		eth:            eth,
		workers:        make([]SimulationWorker, 0, len(sim)),
		workersPerNode: workersPerNode,
	}

	for i := range sim {
		worker := SimulationWorker{
			log:               log.Named("worker").With(zap.Int("worker-id", i)),
			simulationBackend: sim[i],
			simRes:            simRes,
			cancelCache:       cancelCache,
			backgroundWg:      backgroundWg,
		}
		q.workers = append(q.workers, worker)
	}
	return q
}

func (q *SimQueue) Start(ctx context.Context) *sync.WaitGroup {
	process := make([]simqueue.ProcessFunc, 0, len(q.workers)*q.workersPerNode)
	for i := range q.workers {
		if q.workersPerNode > 1 {
			workers := simqueue.MultipleWorkers(q.workers[i].Process, q.workersPerNode, rate.Inf, 1)
			process = append(process, workers...)
		} else {
			process = append(process, q.workers[i].Process)
		}
	}
	blockNumber, err := q.eth.BlockNumber(ctx)
	if err != nil {
		q.log.Warn("Failed to get block number", zap.Error(err))
	} else {
		_ = q.queue.UpdateBlock(blockNumber)
	}

	wg := q.queue.StartProcessLoop(ctx, process)

	wg.Add(1)
	go func() {
		defer wg.Done()

		back := backoff.NewExponentialBackOff()
		back.MaxInterval = 3 * time.Second
		back.MaxElapsedTime = 12 * time.Second

		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := backoff.Retry(func() error {
					blockNumber, err := q.eth.BlockNumber(ctx)
					if err != nil {
						return err
					}
					return q.queue.UpdateBlock(blockNumber)
				}, back)
				if err != nil {
					q.log.Error("Failed to update block number", zap.Error(err))
				}
			}
		}
	}()
	return wg
}

func (q *SimQueue) ScheduleRequest(ctx context.Context, request *SendRequestArgs, highPriority bool) error {
	startAt := time.Now()
	defer func() {
		metrics.RecordRequestAddQueueDuration(time.Since(startAt).Milliseconds())
	}()
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return q.queue.Push(ctx, data, highPriority, uint64(request.Inclusion.DesiredBlock), uint64(request.Inclusion.MaxBlock))
}

type SimulationWorker struct {
	log               *zap.Logger
	simulationBackend SimulationBackend
	simRes            SimulationResult
	cancelCache       *RedisCancellationCache
	backgroundWg      *sync.WaitGroup
}

func (w *SimulationWorker) Process(ctx context.Context, data []byte, info simqueue.QueueItemInfo) (err error) {
	startAt := time.Now()
	defer func() {
		metrics.RecordBundleProcessDuration(time.Since(startAt).Milliseconds())
	}()
	var bundle SendRequestArgs
	err = json.Unmarshal(data, &bundle)
	if err != nil {
		w.log.Error("Failed to unmarshal bundle simulation data", zap.Error(err))
		return err
	}

	var hash common.Hash
	if bundle.Metadata != nil {
		hash = bundle.Metadata.RequestHash
	}
	logger := w.log.With(zap.String("bundle", hash.Hex()))

	// Check if bundle was cancelled
	cancelled, err := w.isBundleCancelled(ctx, &bundle)
	if err != nil {
		// We don't return error here,  because we would consider this error as non-critical as our cancellations are "best effort".
		logger.Error("Failed to check if bundle was cancelled", zap.Error(err))
	}
	if cancelled {
		logger.Info("Bundle is not simulated because it was cancelled")
		return simqueue.ErrProcessUnrecoverable
	}

	w.backgroundWg.Add(1)
	go func() {
		defer w.backgroundWg.Done()
		resCtx, cancel := context.WithTimeout(context.Background(), consumeSimulationTimeout)
		defer cancel()
		err = w.simRes.SimulatedBundle(resCtx, &bundle, info)
		if err != nil {
			w.log.Error("Failed to consume matched share bundle", zap.Error(err))
		}
	}()

	return nil
}

func (w *SimulationWorker) isBundleCancelled(ctx context.Context, bundle *SendRequestArgs) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, simCacheTimeout)
	defer cancel()
	if bundle.Metadata == nil {
		w.log.Error("Bundle has no metadata, skipping cancel check")
		return false, nil
	}
	res, err := w.cancelCache.IsCancelled(ctx, append([]common.Hash{bundle.Metadata.RequestHash}, bundle.Metadata.BodyHashes...))
	if err != nil {
		return false, err
	}
	return res, nil
}
