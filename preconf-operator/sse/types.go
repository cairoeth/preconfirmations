package sse

import (
	"github.com/cairoeth/preconfirmations/preconf-operator/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Event represents a matchmaker event sent from sse subscription
type Event struct {
	Data  *MatchMakerEvent // Will be nil if an error occurred during poll
	Error error
}

// MatchMakerEvent represents the pending transaction hints sent by matchmaker
type MatchMakerEvent struct {
	Hash common.Hash          `json:"hash"`
	Logs []types.Log          `json:"logs,omitempty"`
	Txs  []PendingTransaction `json:"txs,omitempty"`
}

// PendingTransaction represents the hits revealed by the matchmaker about the transaction / bundle
type PendingTransaction struct {
	To               common.Address `json:"to"`
	FunctionSelector [4]byte        `json:"functionSelector,omitempty"`
	CallData         []byte         `json:"callData,omitempty"`
	MevGasPrice      *hexutil.Big   `json:"mevGasPrice,omitempty"`
	GasUsed          *hexutil.Big   `json:"gasUsed,omitempty"`
}
