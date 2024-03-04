package preconshare

import (
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var (
	ErrInvalidHintIntent = errors.New("invalid hint intent")
	ErrNilBundleMetadata = errors.New("bundle metadata is nil")
)

const (
	SendRequestEndpointName    = "preconf_sendRequest"
	ConfirmRequestEndpointName = "preconf_confirmRequest"
	GetRequestEndpointName     = "preconf_getRequest"
)

// HintIntent is a set of hint intents
// its marshalled as an array of strings
type HintIntent uint8

const (
	HintContractAddress HintIntent = 1 << iota
	HintFunctionSelector
	HintLogs
	HintCallData
	HintHash
	HintSpecialLogs
	HintTxHash
	HintNone = 0
)

func (b *HintIntent) SetHint(flag HintIntent) {
	*b = *b | flag
}

func (b *HintIntent) HasHint(flag HintIntent) bool {
	return *b&flag != 0
}

func (b HintIntent) MarshalJSON() ([]byte, error) {
	var arr []string
	if b.HasHint(HintContractAddress) {
		arr = append(arr, "contract_address")
	}
	if b.HasHint(HintFunctionSelector) {
		arr = append(arr, "function_selector")
	}
	if b.HasHint(HintLogs) {
		arr = append(arr, "logs")
	}
	if b.HasHint(HintCallData) {
		arr = append(arr, "calldata")
	}
	if b.HasHint(HintHash) {
		arr = append(arr, "hash")
	}
	if b.HasHint(HintSpecialLogs) {
		arr = append(arr, "special_logs")
	}
	if b.HasHint(HintTxHash) {
		arr = append(arr, "tx_hash")
	}
	return json.Marshal(arr)
}

func (b *HintIntent) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	for _, v := range arr {
		switch v {
		case "contract_address":
			b.SetHint(HintContractAddress)
		case "function_selector":
			b.SetHint(HintFunctionSelector)
		case "logs":
			b.SetHint(HintLogs)
		case "calldata":
			b.SetHint(HintCallData)
		case "hash":
			b.SetHint(HintHash)
		case "special_logs", "default_logs":
			b.SetHint(HintSpecialLogs)
		case "tx_hash":
			b.SetHint(HintTxHash)
		default:
			return ErrInvalidHintIntent
		}
	}
	return nil
}

type Hint struct {
	Hash      common.Hash      `json:"hash"`
	Inclusion RequestInclusion `json:"inclusion"`
	Logs      []CleanLog       `json:"logs"`
	Txs       []TxHint         `json:"txs"`
}

type TxHint struct {
	Hash             *common.Hash    `json:"hash,omitempty"`
	To               *common.Address `json:"to,omitempty"`
	FunctionSelector *hexutil.Bytes  `json:"functionSelector,omitempty"`
	CallData         *hexutil.Bytes  `json:"callData,omitempty"`
}

/////////////////////////
// preconf_sendRequest //
/////////////////////////

type SendRequestArgs struct {
	Version   string           `json:"version"`
	Inclusion RequestInclusion `json:"inclusion"`
	Body      []RequestBody    `json:"body"`
	Privacy   *RequestPrivacy  `json:"privacy,omitempty"`
	Metadata  *RequestMetadata `json:"metadata,omitempty"`
}

type RequestInclusion struct {
	DesiredBlock hexutil.Uint64 `json:"desiredBlock"`
	MaxBlock     hexutil.Uint64 `json:"maxBlock"`
	Tip          hexutil.Uint64 `json:"tip"`
}

type RequestBody struct {
	Tx *hexutil.Bytes `json:"tx,omitempty"`
}

type RequestPrivacy struct {
	Hints     HintIntent `json:"hints,omitempty"`
	Operators []string   `json:"operators,omitempty"`
}

type RequestMetadata struct {
	RequestHash  common.Hash    `json:"requestHash,omitempty"`
	BodyHashes   []common.Hash  `json:"bodyHashes,omitempty"`
	Signer       common.Address `json:"signer,omitempty"`
	OriginID     string         `json:"originId,omitempty"`
	ReceivedAt   hexutil.Uint64 `json:"receivedAt,omitempty"`
	MatchingHash common.Hash    `json:"matchingHash,omitempty"`
	Prematched   bool           `json:"prematched"`
}

type SendRequestResponse struct {
	RequestHash common.Hash    `json:"requestHash"`
	Signature   *hexutil.Bytes `json:"preconfSignature"`
	Block       hexutil.Uint64 `json:"preconfBlock"`
}

////////////////////////////
// preconf_confirmRequest //
////////////////////////////

type ConfirmRequestArgs struct {
	Version   string         `json:"version"`
	Preconf   ConfirmPreconf `json:"preconf"`
	Signature *hexutil.Bytes `json:"signature"`
	Endpoint  string         `json:"endpoint"`
}

type ConfirmPreconf struct {
	Hash  common.Hash    `json:"hash"`
	Block hexutil.Uint64 `json:"block"`
}

type ConfirmRequestResponse struct {
	Valid bool `json:"valid"`
}

////////////////////////////
// preconf_getRequest //
////////////////////////////

type GetRequestArgs struct {
	Version string      `json:"version"`
	Hash    common.Hash `json:"hash"`
}

type GetRequestResponse struct {
	Signature *hexutil.Bytes `json:"signature"`
	Block     hexutil.Uint64 `json:"block"`
	Time      hexutil.Uint64 `json:"time"`
}

type CleanLog struct {
	// address of the contract that generated the event
	Address common.Address `json:"address"`
	// list of topics provided by the contract.
	Topics []common.Hash `json:"topics"`
	// supplied by the contract, usually ABI-encoded
	Data hexutil.Bytes `json:"data"`
}
