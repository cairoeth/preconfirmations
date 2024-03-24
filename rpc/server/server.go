package server

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cairoeth/preconfirmations/rpc/adapters/webfile"
	"github.com/cairoeth/preconfirmations/rpc/application"
	"github.com/cairoeth/preconfirmations/rpc/database"

	"github.com/alicebob/miniredis"
	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var Now = time.Now // used to mock time in tests

var DebugDontSendTx = os.Getenv("DEBUG_DONT_SEND_RAWTX") != ""

// RState Metamask fix helper
var RState *RedisState

type BuilderNameProvider interface {
	BuilderNames() []string
}
type RPCEndPointServer struct {
	server *http.Server
	drain  *http.Server

	drainAddress        string
	drainSeconds        int
	db                  database.Store
	isHealthy           bool
	isHealthyMx         sync.RWMutex
	listenAddress       string
	logger              *zap.Logger
	proxyTimeoutSeconds int
	proxyURL            string
	relaySigningKey     *ecdsa.PrivateKey
	relayURL            string
	startTime           time.Time
	version             string
	builderNameProvider BuilderNameProvider
	chainID             []byte
}

func NewRPCEndPointServer(cfg Configuration) (*RPCEndPointServer, error) {
	var err error
	if DebugDontSendTx {
		cfg.Logger.Info("DEBUG MODE: raw transactions will not be sent out!", zap.String("redisURL", cfg.RedisURL))
	}

	if cfg.RedisURL == "dev" {
		cfg.Logger.Info("Using integrated in-memory Redis instance", zap.String("redisURL", cfg.RedisURL))
		redisServer, err := miniredis.Run()
		if err != nil {
			return nil, err
		}
		cfg.RedisURL = redisServer.Addr()
	}
	// Setup redis connection
	cfg.Logger.Info("Connecting to redis...", zap.String("redisURL", cfg.RedisURL))
	RState, err = NewRedisState(cfg.RedisURL)
	if err != nil {
		return nil, errors.Wrap(err, "Redis init error")
	}
	var builderInfoFetcher application.Fetcher
	if cfg.BuilderInfoSource != "" {
		builderInfoFetcher = webfile.NewFetcher(cfg.BuilderInfoSource)
	}
	bis, err := application.StartBuilderInfoService(context.Background(), builderInfoFetcher, time.Second*time.Duration(cfg.FetchInfoInterval))
	if err != nil {
		return nil, errors.Wrap(err, "BuilderInfoService init error")
	}

	bts, err := fetchNetworkIDBytes(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "fetchNetworkIDBytes error")
	}

	return &RPCEndPointServer{
		db:                  cfg.DB,
		drainAddress:        cfg.DrainAddress,
		drainSeconds:        cfg.DrainSeconds,
		isHealthy:           true,
		listenAddress:       cfg.ListenAddress,
		logger:              cfg.Logger,
		proxyTimeoutSeconds: cfg.ProxyTimeoutSeconds,
		proxyURL:            cfg.ProxyURL,
		relaySigningKey:     cfg.RelaySigningKey,
		relayURL:            cfg.RelayURL,
		startTime:           Now(),
		version:             cfg.Version,
		chainID:             bts,
		builderNameProvider: bis,
	}, nil
}

func fetchNetworkIDBytes(cfg Configuration) ([]byte, error) {
	cl := NewRPCProxyClient(cfg.Logger, cfg.ProxyURL, cfg.ProxyTimeoutSeconds)

	_req := types.NewJSONRPCRequest(1, "net_version", []interface{}{})
	jsonData, err := json.Marshal(_req)
	if err != nil {
		return nil, errors.Wrap(err, "network version request failed")
	}
	httpRes, err := cl.ProxyRequest(jsonData)
	if err != nil {
		return nil, errors.Wrap(err, "cl.NetworkID error")
	}

	resBytes, err := io.ReadAll(httpRes.Body)
	httpRes.Body.Close()
	if err != nil {
		return nil, err
	}
	_res, err := respBytesToJSONRPCResponse(resBytes)
	if err != nil {
		return nil, err
	}

	return _res.Result, nil
}

func (s *RPCEndPointServer) Start() {
	s.logger.Info("Starting rpc endpoint...", zap.String("version", s.version), zap.String("listenAddress", s.listenAddress))

	// Regularly log debug info
	go func() {
		for {
			s.logger.Info("[stats] num-goroutines", zap.Int("count", runtime.NumGoroutine()))
			time.Sleep(10 * time.Second)
		}
	}()

	s.startMainServer()
	s.startDrainServer()

	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, os.Interrupt, syscall.SIGTERM)

	<-notifier

	s.stopDrainServer()
	s.stopMainServer()
}

func (s *RPCEndPointServer) startMainServer() {
	if s.server != nil {
		panic("http server is already running")
	}
	// Handler for root URL (JSON-RPC on POST, public/index.html on GET)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleHTTPRequest)
	mux.HandleFunc("/health", s.handleHealthRequest)
	mux.HandleFunc("/bundle", s.HandleBundleRequest)
	s.server = &http.Server{
		Addr:         s.listenAddress,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("http server failed", zap.Error(err))
		}
	}()
}

func (s *RPCEndPointServer) stopMainServer() {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			s.logger.Error("http server shutdown failed", zap.Error(err))
		}
		s.logger.Info("http server stopped")
		s.server = nil
	}
}

func (s *RPCEndPointServer) startDrainServer() {
	if s.drain != nil {
		panic("drain http server is already running")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleDrain)
	s.drain = &http.Server{
		Addr:    s.drainAddress,
		Handler: mux,
	}
	go func() {
		if err := s.drain.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("drain http server failed", zap.Error(err))
		}
	}()
}

func (s *RPCEndPointServer) stopDrainServer() {
	if s.drain != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.drain.Shutdown(ctx); err != nil {
			s.logger.Error("drain http server shutdown failed", zap.Error(err))
		}
		s.logger.Info("drain http server stopped")
		s.drain = nil
	}
}

func (s *RPCEndPointServer) HandleHTTPRequest(respw http.ResponseWriter, req *http.Request) {
	setCorsHeaders(respw)

	if req.Method == http.MethodGet {
		if strings.Trim(req.URL.Path, "/") == "fast" {
			http.Redirect(respw, req, "https://docs.flashbots.net/flashbots-protect/quick-start#faster-transactions", http.StatusFound)
		} else {
			http.Redirect(respw, req, "https://docs.flashbots.net/flashbots-protect/rpc/quick-start/", http.StatusFound)
		}
		return
	}

	if req.Method == http.MethodOptions {
		respw.WriteHeader(http.StatusOK)
		return
	}

	request := NewRPCRequestHandler(s.logger, &respw, req, s.proxyURL, s.proxyTimeoutSeconds, s.relaySigningKey, s.relayURL, s.db, s.builderNameProvider.BuilderNames(), s.chainID)
	request.process()
}

func (s *RPCEndPointServer) handleDrain(respw http.ResponseWriter, req *http.Request) {
	s.isHealthyMx.Lock()
	if !s.isHealthy {
		s.isHealthyMx.Unlock()
		return
	}

	s.isHealthy = false
	s.logger.Info("Server marked as unhealthy")

	// Let's not hold onto the lock in our sleep
	s.isHealthyMx.Unlock()

	// Give LB enough time to detect us unhealthy
	time.Sleep(
		time.Duration(s.drainSeconds) * time.Second,
	)
}

func (s *RPCEndPointServer) handleHealthRequest(respw http.ResponseWriter, req *http.Request) {
	s.isHealthyMx.RLock()
	defer s.isHealthyMx.RUnlock()
	res := types.HealthResponse{
		Now:       Now(),
		StartTime: s.startTime,
		Version:   s.version,
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		s.logger.Info("[healthCheck] Json error", zap.Error(err))
		respw.WriteHeader(http.StatusInternalServerError)
		return
	}

	respw.Header().Set("Content-Type", "application/json")
	if s.isHealthy {
		respw.WriteHeader(http.StatusOK)
	} else {
		respw.WriteHeader(http.StatusInternalServerError)
	}
	respw.Write(jsonResp)
}

func (s *RPCEndPointServer) HandleBundleRequest(respw http.ResponseWriter, req *http.Request) {
	setCorsHeaders(respw)
	bundleID := req.URL.Query().Get("id")
	if bundleID == "" {
		http.Error(respw, "no bundle id", http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodGet {
		txs, err := RState.GetWhitehatBundleTx(bundleID)
		if err != nil {
			s.logger.Info("[handleBundleRequest] GetWhitehatBundleTx failed", zap.String("bundleID", bundleID), zap.Error(err))
			respw.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := types.BundleResponse{
			BundleID: bundleID,
			RawTxs:   txs,
		}

		jsonResp, err := json.Marshal(res)
		if err != nil {
			s.logger.Info("[handleBundleRequest] Json marshal failed", zap.Error(err))
			respw.WriteHeader(http.StatusInternalServerError)
			return
		}
		respw.Header().Set("Content-Type", "application/json")
		respw.WriteHeader(http.StatusOK)
		respw.Write(jsonResp)
	} else if req.Method == http.MethodDelete {
		RState.DelWhitehatBundleTx(bundleID)
		respw.WriteHeader(http.StatusOK)
	} else {
		respw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func setCorsHeaders(respw http.ResponseWriter) {
	respw.Header().Set("Access-Control-Allow-Origin", "*")
	respw.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type")
}
