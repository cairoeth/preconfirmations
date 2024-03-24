package preconshare

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/ybbus/jsonrpc/v3"
)

type HintBackend interface {
	NotifyHint(ctx context.Context, hint *Hint) error
}

// SimulationBackend is an interface for simulating transactions
// There should be one simulation backend per worker node
type SimulationBackend interface{}

type JSONRPCSimulationBackend struct {
	client jsonrpc.RPCClient
}

func NewJSONRPCSimulationBackend(url string) *JSONRPCSimulationBackend {
	return &JSONRPCSimulationBackend{
		client: jsonrpc.NewClient(url),
		// todo here use optsx
	}
}

type RedisHintBackend struct {
	client     *redis.Client
	pubChannel string
}

func NewRedisHintBackend(redisClient *redis.Client, pubChannel string) *RedisHintBackend {
	return &RedisHintBackend{
		client:     redisClient,
		pubChannel: pubChannel,
	}
}

func (b *RedisHintBackend) NotifyHint(ctx context.Context, hint *Hint) error {
	data, err := json.Marshal(hint)
	if err != nil {
		return err
	}
	// print data
	fmt.Println(string(data))
	return b.client.Publish(ctx, b.pubChannel, data).Err()
}

type JSONRPCBuilder struct {
	url string
}

func (b *JSONRPCBuilder) String() string {
	return b.url
}
