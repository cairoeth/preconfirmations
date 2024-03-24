/* Package server Request represents an incoming client request */
package server

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/cairoeth/preconfirmations/rpc/database"

	"github.com/cairoeth/preconfirmations/preconf-share/preconshare"
	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/zap"
)

type RPCRequest struct {
	logger                     *zap.Logger
	client                     RPCProxyClient
	jsonReq                    *types.JSONRPCRequest
	jsonRes                    *types.JSONRPCResponse
	rawTxHex                   string
	tx                         *ethtypes.Transaction
	txFrom                     string
	relaySigningKey            *ecdsa.PrivateKey
	relayURL                   string
	origin                     string
	referer                    string
	isWhitehatBundleCollection bool
	proxyURL                   string
	ethSendRawTxEntry          *database.EthSendRawTxEntry
	urlParams                  URLParameters
	chainID                    []byte
}

func NewRPCRequest(logger *zap.Logger, client RPCProxyClient, jsonReq *types.JSONRPCRequest, relaySigningKey *ecdsa.PrivateKey, relayURL, origin, referer string, isWhitehatBundleCollection bool, proxyURL string, ethSendRawTxEntry *database.EthSendRawTxEntry, urlParams URLParameters, chainID []byte) *RPCRequest {
	return &RPCRequest{
		logger:                     logger,
		client:                     client,
		jsonReq:                    jsonReq,
		relaySigningKey:            relaySigningKey,
		relayURL:                   relayURL,
		origin:                     origin,
		referer:                    referer,
		isWhitehatBundleCollection: isWhitehatBundleCollection,
		proxyURL:                   proxyURL,
		ethSendRawTxEntry:          ethSendRawTxEntry,
		urlParams:                  urlParams,
		chainID:                    chainID,
	}
}

func (r *RPCRequest) logRequest() {
	if r.jsonReq.Method == "eth_call" && len(r.jsonReq.Params) > 0 {
		p := r.jsonReq.Params[0].(map[string]interface{})
		_to := ""
		_data := ""
		_method := _data
		if p["to"] != nil {
			_to = p["to"].(string)
		}
		if p["data"] != nil {
			_data = p["data"].(string)
		}
		if len(_data) >= 10 {
			_method = _data[:10]
		}
		r.logger.Info("JSON-RPC request", zap.String("method", r.jsonReq.Method), zap.String("paramsTo", _to), zap.String("paramsDataMethod", _method), zap.Int("paramsDataLen", len(_data)), zap.String("origin", r.origin), zap.String("referer", r.referer))
	} else {
		r.logger.Info("JSON-RPC request", zap.String("method", r.jsonReq.Method), zap.String("origin", r.origin), zap.String("referer", r.referer))
	}
}

func (r *RPCRequest) ProcessRequest() *types.JSONRPCResponse {
	r.logRequest()

	switch {
	case r.jsonReq.Method == "eth_sendRawTransaction":
		// r.ethSendRawTxEntry.WhiteHatBundleID = r.whitehatBundleID
		r.handleSendRawTransaction()
	case r.jsonReq.Method == "eth_getTransactionCount" && r.interceptMmEthGetTransactionCount(): // intercept if MM needs to show an error to user
	case r.jsonReq.Method == "eth_call" && r.interceptEthCallToFlashRPCContract(): // intercept if Flashbots isRPC contract
	case r.jsonReq.Method == "net_version":
		r.writeRPCResult(json.RawMessage(r.chainID))
	case r.isWhitehatBundleCollection && r.jsonReq.Method == "eth_getBalance":
		r.writeRPCResult("0x56bc75e2d63100000") // 100 ETH, same as the eth_call SC call above returns
	default:
		if r.isWhitehatBundleCollection && r.jsonReq.Method == "eth_call" {
			r.WhitehatBalanceCheckerRewrite()
		}
		// Proxy the request to a node
		readJSONRPCSuccess := r.proxyRequestRead()
		if !readJSONRPCSuccess {
			r.logger.Info("[ProcessRequest] Proxy to node failed", zap.String("method", r.jsonReq.Method))
			r.writeRPCError("internal server error", types.JSONRPCInternalError)
			return r.jsonRes
		}

		// After proxy, perhaps check backend [MM fix #3 step 2]
		if r.jsonReq.Method == "eth_getTransactionReceipt" {
			requestCompleted := r.checkPostGetTransactionReceipt(r.jsonRes)
			if requestCompleted {
				return r.jsonRes
			}
		}
	}
	return r.jsonRes
}

// Proxies the incoming request to the target URL, and tries to parse JSON-RPC response (and check for specific)
func (r *RPCRequest) proxyRequestRead() (readJSONRpsResponseSuccess bool) {
	timeProxyStart := Now() // for measuring execution time
	body, err := json.Marshal(r.jsonReq)
	if err != nil {
		r.logger.Error("[proxyRequestRead] Failed to marshal request before making proxy request", zap.Error(err))
		return false
	}

	// Proxy request
	proxyResp, err := r.client.ProxyRequest(body)
	if err != nil {
		r.logger.Error("[proxyRequestRead] Failed to make proxy request", zap.Error(err))
		if proxyResp == nil {
			return false
		} else {
			return false
		}
	}

	// Afterwards, check time and result
	timeProxyNeeded := time.Since(timeProxyStart)
	r.logger.Info("[proxyRequestRead] proxied response", zap.Int("secNeeded", int(timeProxyNeeded.Seconds())))

	// Read body
	defer proxyResp.Body.Close()
	proxyRespBody, err := io.ReadAll(proxyResp.Body)
	if err != nil {
		r.logger.Error("[proxyRequestRead] Failed to read proxy request body", zap.Error(err))
		return false
	}

	// Unmarshall JSON-RPC response and check for error inside
	jsonRPCResp := new(types.JSONRPCResponse)
	if err = json.Unmarshal(proxyRespBody, jsonRPCResp); err != nil {
		r.logger.Error("[proxyRequestRead] Failed decoding proxy json-rpc response", zap.Error(err))
		return false
	}
	r.jsonRes = jsonRPCResp
	return true
}

// Check whether to block resending this tx. Send only if (a) not sent before, (b) sent and status=failed, (c) sent, status=unknown and sent at least 5 min ago
func (r *RPCRequest) blockResendingTxToRelay(txHash string) bool {
	timeSent, txWasSentToRelay, err := RState.GetTxSentToRelay(txHash)
	if err != nil {
		r.logger.Error("[blockResendingTxToRelay] Redis:GetTxSentToRelay error", zap.Error(err))
		return false // don't block on redis error
	}

	if !txWasSentToRelay {
		return false // don't block if not sent before
	}

	// was sent before. check status and time
	txStatusAPIResponse, err := GetTxStatus(txHash)
	if err != nil {
		r.logger.Error("[blockResendingTxToRelay] GetTxStatus error", zap.Error(err))
		return false // don't block on redis error
	}

	// Allow sending to relay if tx has failed, or if it's still unknown after a while
	txStatus := txStatusAPIResponse.Status
	if txStatus == types.TxStatusFailed {
		return false // don't block if tx failed
	} else if txStatus == types.TxStatusUnknown && time.Since(timeSent).Minutes() >= 5 {
		return false // don't block if unknown and sent at least 5 min ago
	} else {
		// block tx if pending or already included
		return true
	}
}

// Send tx to relay and finish request (write response)
func (r *RPCRequest) sendTxToRelay() {
	txHash := strings.ToLower(r.tx.Hash().Hex())
	// Check if tx was already forwarded and should be blocked now
	IsBlocked := r.blockResendingTxToRelay(txHash)
	if IsBlocked {
		r.ethSendRawTxEntry.IsBlocked = IsBlocked
		r.logger.Info("[sendTxToRelay] Blocked", zap.String("tx", txHash))
		r.writeRPCResult(txHash)
		return
	}

	r.logger.Info("[sendTxToRelay] sending transaction to relay", zap.String("tx", txHash), zap.String("fromAddress", r.txFrom), zap.String("toAddress", r.tx.To().Hex()))
	r.ethSendRawTxEntry.WasSentToRelay = true

	// mark tx as sent to relay
	err := RState.SetTxSentToRelay(txHash)
	if err != nil {
		r.logger.Error("[sendTxToRelay] Redis:SetTxSentToRelay failed", zap.Error(err))
	}

	minNonce, maxNonce, err := r.GetAddressNonceRange(r.txFrom)
	if err != nil {
		r.logger.Error("[sendTxToRelay] GetAddressNonceRange error", zap.Error(err))
	} else {
		if r.tx.Nonce() < minNonce || r.tx.Nonce() > maxNonce+1 {
			r.logger.Info("[sendTxToRelay] invalid nonce", zap.String("tx", txHash), zap.String("txFrom", r.txFrom), zap.Uint64("minNonce", minNonce), zap.Uint64("maxNonce", maxNonce+1), zap.Uint64("txNonce", r.tx.Nonce()))
			r.writeRPCError("invalid nonce", types.JSONRPCInternalError)
			return
		}
	}

	go RState.SetSenderMaxNonce(r.txFrom, r.tx.Nonce())

	// only allow large transactions to certain addresses - default max tx size is 128KB
	// https://github.com/ethereum/go-ethereum/blob/master/core/tx_pool.go#L53
	if r.tx.Size() > 131072 {
		if r.tx.To() == nil {
			r.logger.Error("[sendTxToRelay] large tx not allowed to target null", zap.String("tx", txHash))
			r.writeRPCError("invalid target for large tx", types.JSONRPCInternalError)
			return
		} else if _, found := allowedLargeTxTargets[strings.ToLower(r.tx.To().Hex())]; !found {
			r.logger.Error("[sendTxToRelay] large tx not allowed to target", zap.String("tx", txHash), zap.String("target", r.tx.To().Hex()))
			r.writeRPCError("invalid target for large tx", types.JSONRPCInternalError)
			return
		}
		r.logger.Info("sendTxToRelay] allowed large tx", zap.String("tx", txHash), zap.String("target", r.tx.To().Hex()))
	}

	// remember this tx based on from+nonce (for cancel-tx)
	err = RState.SetTxHashForSenderAndNonce(r.txFrom, r.tx.Nonce(), txHash)
	if err != nil {
		r.logger.Error("[sendTxToRelay] Redis:SetTxHashForSenderAndNonce failed", zap.Error(err))
	}

	// err = RState.SetLastPrivTxHashOfAccount(r.txFrom, txHash)
	// if err != nil {
	// 	r.Error("[sendTxToRelay] redis:SetLastTxHashOfAccount failed: %v", err)
	// }

	if DebugDontSendTx {
		r.logger.Info("[sendTxToRelay] Faked sending tx to relay, did nothing", zap.String("tx", txHash))
		r.writeRPCResult(txHash)
		return
	}

	sendPrivateTxArgs := types.SendPrivateTxRequestWithPreferences{}
	sendPrivateTxArgs.Tx = r.rawTxHex
	sendPrivateTxArgs.Preferences = &r.urlParams.pref
	if r.urlParams.fast {
		if len(sendPrivateTxArgs.Preferences.Validity.Refund) == 0 {
			addr, err := GetSenderAddressFromTx(r.tx)
			if err != nil {
				r.logger.Error("[sendTxToRelay] GetSenderAddressFromTx failed", zap.Error(err))
				r.writeRPCError(err.Error(), types.JSONRPCInternalError)
				return
			}
			sendPrivateTxArgs.Preferences.Validity.Refund = []types.RefundConfig{
				{
					Address: addr,
					Percent: 50,
				},
			}
		}
	}

	ethBackend, err := ethclient.Dial(r.proxyURL)
	if err != nil {
		r.logger.Fatal("Failed to connect to ethBackend endpoint", zap.Error(err))
	}

	blockNumber, err := ethBackend.BlockNumber(context.Background())
	if err != nil {
		r.logger.Fatal("Failed to connect to get block number", zap.Error(err))
	}

	// Sending transactions to preconf-share node
	rpcClient := jsonrpc.NewClient(r.relayURL)

	txBytes := common.FromHex(r.rawTxHex)

	request := preconshare.SendRequestArgs{
		Version: "v0.1",
		Inclusion: preconshare.RequestInclusion{
			DesiredBlock: hexutil.Uint64(blockNumber),
			MaxBlock:     hexutil.Uint64(blockNumber + 5),
		},
		Body: []preconshare.RequestBody{{Tx: (*hexutil.Bytes)(&txBytes)}},
		Privacy: &preconshare.RequestPrivacy{
			Hints:     preconshare.HintHash,
			Operators: nil,
		},
	}

	var result preconshare.SendRequestResponse

	err = rpcClient.CallFor(context.Background(), &result, "preconf_sendRequest", []*preconshare.SendRequestArgs{&request})
	if err != nil {
		r.logger.Error("[sendTxToRelay] Relay call failed", zap.Error(err))
		r.writeRPCError(err.Error(), types.JSONRPCInternalError)
		return
	}

	r.writeRPCResult(txHash)
	r.logger.Info("[sendTxToRelay] Sent and received preconfirmation", zap.String("tx", txHash), zap.Uint64("block", uint64(result.Block)))
}

func (r *RPCRequest) GetAddressNonceRange(address string) (minNonce, maxNonce uint64, err error) {
	// Get minimum nonce by asking the eth node for the current transaction count
	_req := types.NewJSONRPCRequest(1, "eth_getTransactionCount", []interface{}{r.txFrom, "latest"})
	jsonData, err := json.Marshal(_req)
	if err != nil {
		r.logger.Error("[GetAddressNonceRange] eth_getTransactionCount marshal failed", zap.Error(err))
		return 0, 0, err
	}
	httpRes, err := r.client.ProxyRequest(jsonData)
	if err != nil {
		r.logger.Error("[GetAddressNonceRange] eth_getTransactionCount proxy request failed", zap.Error(err))
		return 0, 0, err
	}

	resBytes, err := io.ReadAll(httpRes.Body)
	httpRes.Body.Close()
	if err != nil {
		r.logger.Error("[GetAddressNonceRange] eth_getTransactionCount read response failed", zap.Error(err))
		return 0, 0, err
	}
	_res, err := respBytesToJSONRPCResponse(resBytes)
	if err != nil {
		r.logger.Error("[GetAddressNonceRange] eth_getTransactionCount parsing response failed", zap.Error(err))
		return 0, 0, err
	}
	_userNonceStr := ""
	err = json.Unmarshal(_res.Result, &_userNonceStr)
	if err != nil {
		r.logger.Error("[GetAddressNonceRange] eth_getTransactionCount unmarshall failed", zap.Error(err))
		r.writeRPCError("internal server error", types.JSONRPCInternalError)
		return
	}
	_userNonceStr = strings.Replace(_userNonceStr, "0x", "", 1)
	_userNonceBigInt := new(big.Int)
	_userNonceBigInt.SetString(_userNonceStr, 16)
	minNonce = _userNonceBigInt.Uint64()

	// Get maximum nonce by looking at redis, which has current pending transactions
	_redisMaxNonce, _, _ := RState.GetSenderMaxNonce(r.txFrom)
	maxNonce = Max(minNonce, _redisMaxNonce)
	return minNonce, maxNonce, nil
}

func (r *RPCRequest) WhitehatBalanceCheckerRewrite() {
	var err error

	if len(r.jsonReq.Params) == 0 {
		return
	}

	// Ensure param is of type map
	t := reflect.TypeOf(r.jsonReq.Params[0])
	if t.Kind() != reflect.Map {
		return
	}

	p := r.jsonReq.Params[0].(map[string]interface{})
	if to := p["to"]; to == "0xb1f8e55c7f64d203c1400b9d8555d050f94adf39" {
		r.jsonReq.Params[0].(map[string]interface{})["to"] = "0x268F7Cd7A396BCE178f0937095772C7fb83a9104"
		if err != nil {
			r.logger.Error("[WhitehatBalanceCheckerRewrite] isWhitehatBundleCollection json marshal failed:", zap.Error(err))
		} else {
			r.logger.Info("[WhitehatBalanceCheckerRewrite] BalanceChecker contract was rewritten to new version")
		}
	}
}

func (r *RPCRequest) writeRPCError(msg string, errCode int) {
	if r.jsonReq.Method == "eth_sendRawTransaction" {
		r.ethSendRawTxEntry.Error = msg
		r.ethSendRawTxEntry.ErrorCode = errCode
	}
	r.jsonRes = &types.JSONRPCResponse{
		ID:      r.jsonReq.ID,
		Version: "2.0",
		Error: &types.JSONRPCError{
			Code:    errCode,
			Message: msg,
		},
	}
}

func (r *RPCRequest) writeRPCResult(result interface{}) {
	resBytes, err := json.Marshal(result)
	if err != nil {
		r.logger.Error("[writeRPCResult] writeRPCResult error marshalling", zap.Error(err))
		r.writeRPCError("internal server error", types.JSONRPCInternalError)
		return
	}
	r.jsonRes = &types.JSONRPCResponse{
		ID:      r.jsonReq.ID,
		Version: "2.0",
		Result:  resBytes,
	}
}
