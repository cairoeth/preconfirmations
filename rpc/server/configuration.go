package server

import (
	"crypto/ecdsa"

	"github.com/cairoeth/preconfirmations/rpc/database"
	"go.uber.org/zap"
)

type Configuration struct {
	DB                  database.Store
	DrainAddress        string
	DrainSeconds        int
	ListenAddress       string
	Logger              *zap.Logger
	ProxyTimeoutSeconds int
	ProxyUrl            string
	RedisUrl            string
	RelaySigningKey     *ecdsa.PrivateKey
	RelayUrl            string
	Version             string
	BuilderInfoSource   string
	FetchInfoInterval   int
}
