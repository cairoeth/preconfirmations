package main

import (
	"crypto/ecdsa"
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/cairoeth/preconfirmations-avs/rpc/database"
	"github.com/cairoeth/preconfirmations-avs/rpc/server"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	version = "dev" // is set during build process

	// defaults
	defaultDebug                    = os.Getenv("DEBUG") == "1"
	defaultLogJSON                  = os.Getenv("LOG_JSON") == "1"
	defaultListenAddress            = "127.0.0.1:9000"
	defaultDrainAddress             = "127.0.0.1:9001"
	defaultDrainSeconds             = 60
	defaultProxyUrl                 = "http://127.0.0.1:8545"
	defaultProxyTimeoutSeconds      = 10
	defaultRelayUrl                 = "https://relay.flashbots.net"
	defaultRedisUrl                 = "localhost:6379"
	defaultServiceName              = os.Getenv("SERVICE_NAME")
	defaultFetchInfoIntervalSeconds = 600

	// cli flags
	versionPtr           = flag.Bool("version", false, "just print the program version")
	listenAddress        = flag.String("listen", getEnvAsStrOrDefault("LISTEN_ADDR", defaultListenAddress), "Listen address")
	drainAddress         = flag.String("drain", getEnvAsStrOrDefault("DRAIN_ADDR", defaultDrainAddress), "Drain address")
	drainSeconds         = flag.Int("drainSeconds", getEnvAsIntOrDefault("DRAIN_SECONDS", defaultDrainSeconds), "seconds to wait for graceful shutdown")
	fetchIntervalSeconds = flag.Int("fetchIntervalSeconds", getEnvAsIntOrDefault("FETCH_INFO_INTERVAL_SECONDS", defaultFetchInfoIntervalSeconds), "seconds between builder info fetches")
	builderInfoSource    = flag.String("builderInfoSource", getEnvAsStrOrDefault("BUILDER_INFO_SOURCE", ""), "URL for json source of actual builder info")
	proxyUrl             = flag.String("proxy", getEnvAsStrOrDefault("PROXY_URL", defaultProxyUrl), "URL for default JSON-RPC proxy target (eth node, Infura, etc.)")
	proxyTimeoutSeconds  = flag.Int("proxyTimeoutSeconds", getEnvAsIntOrDefault("PROXY_TIMEOUT_SECONDS", defaultProxyTimeoutSeconds), "proxy client timeout in seconds")
	redisUrl             = flag.String("redis", getEnvAsStrOrDefault("REDIS_URL", defaultRedisUrl), "URL for Redis (use 'dev' to use integrated in-memory redis)")
	relayUrl             = flag.String("relayUrl", getEnvAsStrOrDefault("RELAY_URL", defaultRelayUrl), "URL for relay")
	relaySigningKey      = flag.String("signingKey", os.Getenv("RELAY_SIGNING_KEY"), "Signing key for relay requests")
	psqlDsn              = flag.String("psql", os.Getenv("POSTGRES_DSN"), "Postgres DSN")
	debugPtr             = flag.Bool("debug", defaultDebug, "print debug output")
	logJSONPtr           = flag.Bool("logJSON", defaultLogJSON, "log in JSON")
	serviceName          = flag.String("serviceName", defaultServiceName, "name of the service which will be used in the logs")
)

func main() {
	var key *ecdsa.PrivateKey
	var err error

	flag.Parse()
	// logFormat := log.TerminalFormat(true)
	// if *logJSONPtr {
	// 	logFormat = log.JSONFormat()
	// }

	// logLevel := log.LvlInfo
	// if *debugPtr {
	// 	logLevel = log.LvlDebug
	// }

	// log.Root().SetHandler(log.LvlFilterHandler(logLevel, log.StreamHandler(os.Stderr, logFormat)))
	// logger := log.New()
	// if *serviceName != "" {
	// 	logger = logger.New(log.Ctx{"service": *serviceName})
	// }
	// // Perhaps print only the version
	// if *versionPtr {
	// 	logger.Info("rpc-endpoint", "version", version)
	// 	return
	// }

	logger, _ := zap.NewDevelopment()
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
	s, err := server.NewRpcEndPointServer(server.Configuration{
		DB:                  db,
		DrainAddress:        *drainAddress,
		DrainSeconds:        *drainSeconds,
		ListenAddress:       *listenAddress,
		Logger:              logger,
		ProxyTimeoutSeconds: *proxyTimeoutSeconds,
		ProxyUrl:            *proxyUrl,
		RedisUrl:            *redisUrl,
		RelaySigningKey:     key,
		RelayUrl:            *relayUrl,
		Version:             version,
		BuilderInfoSource:   *builderInfoSource,
		FetchInfoInterval:   *fetchIntervalSeconds,
	})
	if err != nil {
		logger.Error("Server init error", zap.Error(err))
	}
	logger.Info("Starting rpc-endpoint...", zap.String("relayUrl", *relayUrl), zap.String("proxyUrl", *proxyUrl))
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