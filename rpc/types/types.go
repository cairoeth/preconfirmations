// Package types contains the types used in the JSON-RPC requests and responses.
package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/metachris/flashbotsrpc"
)

// As per JSON-RPC 2.0 Specification
// https://www.jsonrpc.org/specification#error_object
const (
	JSONRPCParseError     = -32700
	JSONRPCInvalidRequest = -32600
	JSONRPCMethodNotFound = -32601
	JSONRPCInvalidParams  = -32602
	JSONRPCInternalError  = -32603
)

type JSONRPCRequest struct {
	ID      interface{}   `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Version string        `json:"jsonrpc,omitempty"`
}

func NewJSONRPCRequest(id interface{}, method string, params []interface{}) *JSONRPCRequest {
	return &JSONRPCRequest{
		ID:      id,
		Method:  method,
		Params:  params,
		Version: "2.0",
	}
}

func NewJSONRPCRequest1(id interface{}, method string, param interface{}) *JSONRPCRequest {
	return NewJSONRPCRequest(id, method, []interface{}{param})
}

type JSONRPCResponse struct {
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	Version string          `json:"jsonrpc"`
}

// JSONRPCError https://www.jsonrpc.org/specification#error_object
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err JSONRPCError) Error() string {
	return fmt.Sprintf("Error %d (%s)", err.Code, err.Message)
}

func NewJSONRPCResponse(id interface{}, result json.RawMessage) *JSONRPCResponse {
	return &JSONRPCResponse{
		ID:      id,
		Result:  result,
		Version: "2.0",
	}
}

type GetBundleStatusByTransactionHashResponse struct {
	TxHash            string `json:"txHash"`            // "0x0aeb9c61b342f7fc94a10d41c5d30a049a9cfa9ab764c6dd02204a19960ee567"
	Status            string `json:"status"`            // "FAILED_BUNDLE"
	Message           string `json:"message"`           // "Expired - The base fee was to low to execute this transaction, please try again"
	Error             string `json:"error"`             // "max fee per gas less than block base fee"
	BlocksCount       int    `json:"blocksCount"`       // 2
	ReceivedTimestamp int    `json:"receivedTimestamp"` // 1634568851003
	StatusTimestamp   int    `json:"statusTimestamp"`   // 1634568873862
}

type HealthResponse struct {
	Now       time.Time `json:"time"`
	StartTime time.Time `json:"startTime"`
	Version   string    `json:"version"`
}

type TransactionReceipt struct {
	TransactionHash string
	Status          string
}

type PrivateTxStatus string

var (
	TxStatusUnknown  PrivateTxStatus = "UNKNOWN"
	TxStatusPending  PrivateTxStatus = "PENDING"
	TxStatusIncluded PrivateTxStatus = "INCLUDED"
	TxStatusFailed   PrivateTxStatus = "FAILED"
)

type PrivateTxAPIResponse struct {
	Status         PrivateTxStatus `json:"status"`
	Hash           string          `json:"hash"`
	MaxBlockNumber int             `json:"maxBlockNumber"`
}

type RelayErrorResponse struct {
	Error string `json:"error"`
}

type BundleResponse struct {
	BundleID string   `json:"bundleID"`
	RawTxs   []string `json:"rawTxs"`
}

type SendPrivateTxRequestWithPreferences struct {
	flashbotsrpc.FlashbotsSendPrivateTransactionRequest
	Preferences *PrivateTxPreferences `json:"preferences,omitempty"`
}

type TxPrivacyPreferences struct {
	Hints    []string `json:"hints"`
	Builders []string `json:"builders"`
}

type TxValidityPreferences struct {
	Refund []RefundConfig `json:"refund,omitempty"`
}

type RefundConfig struct {
	Address common.Address `json:"address"`
	Percent int            `json:"percent"`
}

type PrivateTxPreferences struct {
	Privacy  TxPrivacyPreferences  `json:"privacy"`
	Validity TxValidityPreferences `json:"validity"`
}
