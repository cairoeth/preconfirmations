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
	ProxyURL            string
	RedisURL            string
	RelaySigningKey     *ecdsa.PrivateKey
	RelayURL            string
	Version             string
	BuilderInfoSource   string
	FetchInfoInterval   int
}
