package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/cairoeth/preconfirmations/preconf-share/jsonrpcserver"
	"github.com/cairoeth/preconfirmations/preconf-share/preconshare"
	"github.com/cairoeth/preconfirmations/preconf-share/simqueue"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/flashbots/go-utils/cli"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	version = "dev" // is set during build process

	// Simqueue is configured using its own env variables, see `simqueue` package.

	// Default values
	defaultDebug               = os.Getenv("DEBUG") == "1"
	defaultLogProd             = os.Getenv("LOG_PROD") == "1"
	defaultLogService          = os.Getenv("LOG_SERVICE")
	defaultPort                = cli.GetEnv("PORT", "8080")
	defaultMetricsPort         = cli.GetEnv("METRICS_PORT", "8088")
	defaultChannelName         = cli.GetEnv("REDIS_CHANNEL_NAME", "hints")
	defaultRedisEndpoint       = cli.GetEnv("REDIS_ENDPOINT", "redis://localhost:6379")
	defaultSimulationsEndpoint = cli.GetEnv("SIMULATION_ENDPOINTS", "http://127.0.0.1:8545")
	defaultWorkersPerNode      = cli.GetEnv("WORKERS_PER_SIM_ENDPOINT", "2")
	defaultPostgresDSN         = cli.GetEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	defaultEthEndpoint         = cli.GetEnv("ETH_ENDPOINT", "http://127.0.0.1:8545")

	// Flags
	debugPtr          = flag.Bool("debug", defaultDebug, "print debug output")
	logProdPtr        = flag.Bool("log-prod", defaultLogProd, "log in production mode (json)")
	logServicePtr     = flag.String("log-service", defaultLogService, "'service' tag to logs")
	portPtr           = flag.String("port", defaultPort, "port to listen on")
	channelPtr        = flag.String("channel", defaultChannelName, "redis pub/sub channel name string")
	redisPtr          = flag.String("redis", defaultRedisEndpoint, "redis url string")
	simEndpointPtr    = flag.String("sim-endpoint", defaultSimulationsEndpoint, "simulation endpoints (comma separated)")
	workersPerNodePtr = flag.String("workers-per-node", defaultWorkersPerNode, "number of workers per simulation node")
	postgresDSNPtr    = flag.String("postgres-dsn", defaultPostgresDSN, "postgres dsn")
	ethPtr            = flag.String("eth", defaultEthEndpoint, "eth endpoint")
)

func main() {
	flag.Parse()

	logger, _ := zap.NewDevelopment()
	if *logProdPtr {
		atom := zap.NewAtomicLevel()
		if *debugPtr {
			atom.SetLevel(zap.DebugLevel)
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		logger = zap.New(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		))
	}
	defer func() { _ = logger.Sync() }()
	if *logServicePtr != "" {
		logger = logger.With(zap.String("service", *logServicePtr))
	}

	ctx, ctxCancel := context.WithCancel(context.Background())

	logger.Info("Starting preconfirmations-avs", zap.String("version", version))

	redisOpts, err := redis.ParseURL(*redisPtr)
	if err != nil {
		logger.Fatal("Failed to parse redis url", zap.Error(err))
	}
	redisClient := redis.NewClient(redisOpts)

	var simBackends []preconshare.SimulationBackend //nolint:prealloc
	for _, simEndpoint := range strings.Split(*simEndpointPtr, ",") {
		simBackend := preconshare.NewJSONRPCSimulationBackend(simEndpoint)
		simBackends = append(simBackends, simBackend)
	}

	hintBackend := preconshare.NewRedisHintBackend(redisClient, *channelPtr)
	if err != nil {
		logger.Fatal("Failed to create redis hint backend", zap.Error(err))
	}

	ethBackend, err := ethclient.Dial(*ethPtr)
	if err != nil {
		logger.Fatal("Failed to connect to ethBackend endpoint", zap.Error(err))
	}

	dbBackend, err := preconshare.NewDBBackend(*postgresDSNPtr)
	if err != nil {
		logger.Fatal("Failed to create postgres backend", zap.Error(err))
	}

	simResultBackend := preconshare.NewSimulationResultBackend(logger, hintBackend, ethBackend, dbBackend)

	redisQueue := simqueue.NewRedisQueue(logger, redisClient, "node")
	redisQueueConfig, err := simqueue.ConfigFromEnv()
	if err != nil {
		logger.Fatal("Failed to load redis queue config", zap.Error(err))
	}
	redisQueue.Config = redisQueueConfig

	// keep track of cancelled bundles for a 30-block window
	cancelCache := preconshare.NewRedisCancellationCache(redisClient, 30*12*time.Second, "node-cancel")

	var workersPerNode int
	if _, err := fmt.Sscanf(*workersPerNodePtr, "%d", &workersPerNode); err != nil {
		logger.Fatal("Failed to parse workers per node", zap.Error(err))
	}
	if workersPerNode < 1 {
		logger.Fatal("Workers per node must be greater than 0")
	}
	backgroundWg := &sync.WaitGroup{}
	simQueue := preconshare.NewQueue(logger, redisQueue, ethBackend, simBackends, simResultBackend, workersPerNode, backgroundWg, cancelCache)
	queueWg := simQueue.Start(ctx)
	// chain id
	chainID, err := ethBackend.ChainID(ctx)
	if err != nil {
		logger.Fatal("Failed to get chain id", zap.Error(err))
	}
	signer := types.LatestSignerForChainID(chainID)

	cachingEthBackend := preconshare.NewCachingEthClient(ethBackend)

	api := preconshare.NewAPI(logger, simQueue, dbBackend, cachingEthBackend, signer, simBackends, cancelCache, time.Millisecond*60)

	jsonRPCServer, err := jsonrpcserver.NewHandler(jsonrpcserver.Methods{
		preconshare.SendRequestEndpointName:    api.SendRequest,
		preconshare.ConfirmRequestEndpointName: api.ConfirmRequest,
		preconshare.GetRequestEndpointName:     api.GetRequest,
	})
	if err != nil {
		logger.Fatal("Failed to create jsonrpc server", zap.Error(err))
	}

	http.Handle("/", jsonRPCServer)
	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", *portPtr),
		ReadHeaderTimeout: 5 * time.Second,
	}

	metricsMux := http.NewServeMux()
	metricsMux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	go func() {
		metricsMux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		metricsMux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		metricsMux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		metricsMux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		metricsMux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

		metricsServer := &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%s", defaultMetricsPort),
			ReadHeaderTimeout: 5 * time.Second,
			Handler:           metricsMux,
		}

		err := metricsServer.ListenAndServe()
		if err != nil {
			logger.Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	connectionsClosed := make(chan struct{})
	go func() {
		notifier := make(chan os.Signal, 1)
		signal.Notify(notifier, os.Interrupt, syscall.SIGTERM)
		<-notifier
		logger.Info("Shutting down...")
		ctxCancel()
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown server", zap.Error(err))
		}
		close(connectionsClosed)
	}()

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("ListenAndServe: ", zap.Error(err))
	}

	<-ctx.Done()
	<-connectionsClosed
	// wait for queue to finish processing
	queueWg.Wait()
	backgroundWg.Wait()
}
