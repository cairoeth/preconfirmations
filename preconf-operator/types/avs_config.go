// Package types contains the types used in the AVS JSON-RPC communication.
package types

import (
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// Log - Custom type because of hex string to bytes decoding error while using default geth.Log
type Log struct {
	Address common.Address `json:"address"`
	Topics  []common.Hash  `json:"topics"`
	Data    []byte         `json:"data,omitempty"` // Could be replaced with geth.hexutil type
}

// UnmarshalJSON unmarshals JSON data into a Log.
func (l *Log) UnmarshalJSON(data []byte) error {
	var temp struct {
		Address common.Address `json:"address"`
		Topics  []common.Hash  `json:"topics"`
		Data    string         `json:"data,omitempty"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	l.Topics = temp.Topics
	l.Address = temp.Address

	if temp.Data != "" {
		temp.Data = temp.Data[2:]

		decoded, err := hex.DecodeString(temp.Data)
		if err != nil {
			return err
		}

		l.Data = decoded
	} else {
		l.Data = nil
	}

	return nil
}

type NodeConfig struct {
	// used to set the logger level (true = info, false = debug)
	Production                    bool   `yaml:"production"`
	OperatorAddress               string `yaml:"operator_address"`
	OperatorStateRetrieverAddress string `yaml:"operator_state_retriever_address"`
	AVSRegistryCoordinatorAddress string `yaml:"avs_registry_coordinator_address"`
	TokenStrategyAddr             string `yaml:"token_strategy_addr"`
	EthRPCURL                     string `yaml:"eth_rpc_url"`
	BlsPrivateKeyStorePath        string `yaml:"bls_private_key_store_path"`
	EcdsaPrivateKeyStorePath      string `yaml:"ecdsa_private_key_store_path"`
	AggregatorServerIPPortAddress string `yaml:"aggregator_server_ip_port_address"`
	RegisterOperatorOnStartup     bool   `yaml:"register_operator_on_startup"`
	EigenMetricsIPPortAddress     string `yaml:"eigen_metrics_ip_port_address"`
	EnableMetrics                 bool   `yaml:"enable_metrics"`
	NodeAPIIPPortAddress          string `yaml:"node_api_ip_port_address"`
	EnableNodeAPI                 bool   `yaml:"enable_node_api"`
}
