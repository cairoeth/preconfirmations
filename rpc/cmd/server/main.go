package main

import (
	"crypto/ecdsa"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/cairoeth/preconfirmations/rpc/database"
	"github.com/cairoeth/preconfirmations/rpc/server"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	version = "dev" // is set during build process

	// defaults
	defaultDebug                    = os.Getenv("DEBUG") == "1"
	defaultListenAddress            = "127.0.0.1:9000"
	defaultDrainAddress             = "127.0.0.1:9001"
	defaultDrainSeconds             = 60
	defaultProxyURL                 = "http://127.0.0.1:8545"
	defaultProxyTimeoutSeconds      = 10
	defaultRelayURL                 = "http://localhost:8080"
	defaultRedisURL                 = "localhost:6379"
	defaultFetchInfoIntervalSeconds = 600

	// cli flags
	listenAddress        = flag.String("listen", getEnvAsStrOrDefault("LISTEN_ADDR", defaultListenAddress), "Listen address")
	drainAddress         = flag.String("drain", getEnvAsStrOrDefault("DRAIN_ADDR", defaultDrainAddress), "Drain address")
	drainSeconds         = flag.Int("drainSeconds", getEnvAsIntOrDefault("DRAIN_SECONDS", defaultDrainSeconds), "seconds to wait for graceful shutdown")
	fetchIntervalSeconds = flag.Int("fetchIntervalSeconds", getEnvAsIntOrDefault("FETCH_INFO_INTERVAL_SECONDS", defaultFetchInfoIntervalSeconds), "seconds between builder info fetches")
	builderInfoSource    = flag.String("builderInfoSource", getEnvAsStrOrDefault("BUILDER_INFO_SOURCE", ""), "URL for json source of actual builder info")
	proxyURL             = flag.String("proxy", getEnvAsStrOrDefault("PROXY_URL", defaultProxyURL), "URL for default JSON-RPC proxy target (eth node, Infura, etc.)")
	proxyTimeoutSeconds  = flag.Int("proxyTimeoutSeconds", getEnvAsIntOrDefault("PROXY_TIMEOUT_SECONDS", defaultProxyTimeoutSeconds), "proxy client timeout in seconds")
	redisURL             = flag.String("redis", getEnvAsStrOrDefault("REDIS_URL", defaultRedisURL), "URL for Redis (use 'dev' to use integrated in-memory redis)")
	relayURL             = flag.String("relayURL", getEnvAsStrOrDefault("RELAY_URL", defaultRelayURL), "URL for preconf rpc")
	relaySigningKey      = flag.String("signingKey", os.Getenv("RELAY_SIGNING_KEY"), "Signing key for relay requests")
	psqlDsn              = flag.String("psql", os.Getenv("POSTGRES_DSN"), "Postgres DSN")
	debugPtr             = flag.Bool("debug", defaultDebug, "print debug output")
)

func main() {
	var key *ecdsa.PrivateKey
	var err error

	flag.Parse()

	atom := zap.NewAtomicLevel()
	if *debugPtr {
		atom.SetLevel(zap.DebugLevel)
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	logger.Info("Init rpc-endpoint", zap.String("version", version))

	if *relaySigningKey == "" {
		logger.Error("Cannot use the relay without a signing key.")
	}

	pkHex := strings.Replace(*relaySigningKey, "0x", "", 1)
	if pkHex == "dev" {
		logger.Info("Creating a new dev signing key...")
		key, err = crypto.GenerateKey()
	} else {
		key, err = crypto.HexToECDSA(pkHex)
	}

	if err != nil {
		logger.Error("Error with relay signing key", zap.Error(err))
	}

	// Setup database
	var db database.Store
	if *psqlDsn == "" {
		db = database.NewMockStore()
	} else {
		db = database.NewPostgresStore(*psqlDsn)
	}
	// Start the endpoint
	s, err := server.NewRPCEndPointServer(server.Configuration{
		DB:                  db,
		DrainAddress:        *drainAddress,
		DrainSeconds:        *drainSeconds,
		ListenAddress:       *listenAddress,
		Logger:              logger,
		ProxyTimeoutSeconds: *proxyTimeoutSeconds,
		ProxyURL:            *proxyURL,
		RedisURL:            *redisURL,
		RelaySigningKey:     key,
		RelayURL:            *relayURL,
		Version:             version,
		BuilderInfoSource:   *builderInfoSource,
		FetchInfoInterval:   *fetchIntervalSeconds,
	})
	if err != nil {
		logger.Error("Server init error", zap.Error(err))
	}
	logger.Info("Starting rpc-endpoint...", zap.String("relayURL", *relayURL), zap.String("proxyURL", *proxyURL))
	s.Start()
}

func getEnvAsStrOrDefault(key, defaultValue string) string {
	ret := os.Getenv(key)
	if ret == "" {
		ret = defaultValue
	}
	return ret
}

func getEnvAsIntOrDefault(name string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
