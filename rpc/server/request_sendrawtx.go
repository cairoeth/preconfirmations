package server

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/ethereum/go-ethereum/common/hexutil"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

const (
	scMethodBytes = 4 // first 4 byte of data field
)

func (r *RPCRequest) handleSendRawTransaction() {
	var err error

	// JSON-RPC sanity checks
	if len(r.jsonReq.Params) < 1 {
		r.logger.Info("[sendRawTransaction] No params for eth_sendRawTransaction")
		r.writeRPCError("empty params for eth_sendRawTransaction", types.JSONRPCInvalidParams)
		return
	}

	if r.jsonReq.Params[0] == nil {
		r.logger.Info("[sendRawTransaction]  Nil param for eth_sendRawTransaction")
		r.writeRPCError("nil params for eth_sendRawTransaction", types.JSONRPCInvalidParams)
	}

	r.rawTxHex = r.jsonReq.Params[0].(string)
	if len(r.rawTxHex) < 2 {
		r.logger.Error("[sendRawTransaction] Invalid raw transaction (wrong length)")
		r.writeRPCError("invalid raw transaction param (wrong length)", types.JSONRPCInvalidParams)
		return
	}

	// r.logger.Info("[sendRawTransaction] Raw tx value", "tx", r.rawTxHex, "txHash", r.tx.Hash())
	r.ethSendRawTxEntry.TxRaw = r.rawTxHex
	r.tx, err = GetTx(r.rawTxHex)
	if err != nil {
		r.logger.Info("[sendRawTransaction] Reading transaction object failed", zap.String("tx", r.rawTxHex))
		r.writeRPCError(fmt.Sprintf("reading transaction object failed - rawTx: %s", r.rawTxHex), types.JSONRPCInvalidRequest)
		return
	}
	r.ethSendRawTxEntry.TxHash = r.tx.Hash().String()
	// Get address from tx
	r.txFrom, err = GetSenderFromRawTx(r.tx)
	if err != nil {
		r.logger.Info("[sendRawTransaction] Couldn't get address from rawTx", zap.Error(err))
		r.writeRPCError(fmt.Sprintf("couldn't get address from rawTx: %v", err), types.JSONRPCInvalidRequest)
		return
	}

	r.logger.Info("[sendRawTransaction] sending raw transaction", zap.String("tx", r.rawTxHex), zap.String("txHash", r.tx.Hash().Hex()), zap.String("fromAddress", r.txFrom), zap.String("toAddress", AddressPtrToStr(r.tx.To())), zap.Uint64("txNonce", r.tx.Nonce()))
	txFromLower := strings.ToLower(r.txFrom)

	// store tx info to ethSendRawTxEntries which will be stored in db for data analytics reason
	r.ethSendRawTxEntry.TxFrom = r.txFrom
	r.ethSendRawTxEntry.TxTo = AddressPtrToStr(r.tx.To())
	r.ethSendRawTxEntry.TxNonce = int(r.tx.Nonce())

	if len(r.tx.Data()) > 0 {
		r.ethSendRawTxEntry.TxData = hexutil.Encode(r.tx.Data())
	}

	if len(r.tx.Data()) >= scMethodBytes {
		r.ethSendRawTxEntry.TxSmartContractMethod = hexutil.Encode(r.tx.Data()[:scMethodBytes])
	}

	if r.tx.Nonce() >= 1e9 {
		r.logger.Info("[sendRawTransaction] tx rejected - nonce too high", zap.Uint64("txNonce", r.tx.Nonce()), zap.String("txHash", r.tx.Hash().Hex()), zap.String("txFromLower", txFromLower), zap.String("origin", r.origin))
		r.writeRPCError("tx rejected - nonce too high", types.JSONRPCInvalidRequest)
		return
	}

	txHashLower := strings.ToLower(r.tx.Hash().Hex())
	// Check if tx was blocked (eg. "nonce too low")
	retVal, isBlocked, _ := RState.GetBlockedTxHash(txHashLower)
	if isBlocked {
		r.logger.Info("[sendRawTransaction] tx blocked", zap.String("txHash", r.tx.Hash().Hex()), zap.String("retVal", retVal))
		r.writeRPCError(retVal, types.JSONRPCInternalError)
		return
	}

	// Remember sender of the tx, for lookup in getTransactionReceipt to possibly set nonce-fix
	err = RState.SetSenderOfTxHash(txHashLower, txFromLower)
	if err != nil {
		r.logger.Error("[sendRawTransaction] Redis:SetSenderOfTxHash failed", zap.Error(err))
	}
	var txToAddr string
	if r.tx.To() != nil { // to address will be nil for contract creation tx
		txToAddr = r.tx.To().String()
	}
	isOnOfacList := isOnOFACList(r.txFrom) || isOnOFACList(txToAddr)
	r.ethSendRawTxEntry.IsOnOafcList = isOnOfacList
	if isOnOfacList {
		r.logger.Info("[sendRawTransaction] Blocked tx due to ofac sanctioned address", zap.String("txHash", txHashLower), zap.String("txFrom", r.txFrom), zap.String("txTo", txToAddr))
		r.writeRPCError("blocked tx due to ofac sanctioned address", types.JSONRPCInvalidRequest)
		return
	}

	// Check if transaction needs protection
	needsProtection := r.doesTxNeedFrontrunningProtection(r.tx)
	r.ethSendRawTxEntry.NeedsFrontRunningProtection = needsProtection
	// If users specify a bundle ID, cache this transaction
	// if r.isWhitehatBundleCollection {
	// 	r.logger.Info("[WhitehatBundleCollection] Adding tx to bundle", zap.String("whiteHatBundleId", r.whitehatBundleID), zap.String("tx", r.rawTxHex))
	// 	err = RState.AddTxToWhitehatBundle(r.whitehatBundleID, r.rawTxHex)
	// 	if err != nil {
	// 		r.logger.Error("[WhitehatBundleCollection] AddTxToWhitehatBundle failed", zap.Error(err))
	// 		r.writeRPCError("[WhitehatBundleCollection] AddTxToWhitehatBundle failed:", types.JSONRPCInternalError)
	// 		return
	// 	}
	// 	r.writeRPCResult(r.tx.Hash().Hex())
	// 	return
	// }

	// // Check for cancellation-tx
	// if r.tx.To() != nil && len(r.tx.Data()) <= 2 && txFromLower == strings.ToLower(r.tx.To().Hex()) {
	// 	r.ethSendRawTxEntry.IsCancelTx = true
	// 	requestDone := r.handleCancelTx() // returns true if tx was cancelled at the relay and response has been sent to the user
	// 	if requestDone {                  // a cancel-tx to fast endpoint is also sent to mempool
	// 		return
	// 	}

	// 	// It's a cancel-tx for the mempool
	// 	needsProtection = false
	// 	r.logger.Info("[cancel-tx] Sending to mempool", zap.String("txFromLower", txFromLower), zap.Uint64("txNonce", r.tx.Nonce()))
	// }

	needsProtection = true

	if needsProtection {
		r.sendTxToRelay()
		return
	}

	if DebugDontSendTx {
		r.logger.Info("[sendRawTransaction] Faked sending tx to mempool, did nothing")
		r.writeRPCResult(r.tx.Hash().Hex())
		return
	}

	// Proxy to public node now
	readJSONRPCSuccess := r.proxyRequestRead()
	r.ethSendRawTxEntry.WasSentToMempool = true
	// Log after proxying
	if !readJSONRPCSuccess {
		r.logger.Error("[sendRawTransaction] Proxy to mempool failed")
		r.writeRPCError("internal server error", types.JSONRPCInternalError)
		return
	}

	// at the end, save the nonce for further spam protection checks
	go RState.SetSenderMaxNonce(txFromLower, r.tx.Nonce())

	if r.jsonRes.Error != nil {
		r.logger.Info("[sendRawTransaction] Proxied eth_sendRawTransaction to mempool", zap.String("jsonRpcError", r.jsonRes.Error.Message), zap.String("txHash", r.tx.Hash().Hex()))
		r.ethSendRawTxEntry.Error = r.jsonRes.Error.Message
		r.ethSendRawTxEntry.ErrorCode = r.jsonRes.Error.Code
		if r.jsonRes.Error.Message == "nonce too low" {
			RState.SetBlockedTxHash(txHashLower, "nonce too low")
		}
	} else {
		r.logger.Info("[sendRawTransaction] Proxied eth_sendRawTransaction to mempool", zap.String("txHash", r.tx.Hash().Hex()))
	}
}

// Check if a request needs frontrunning protection. There are many transactions that don't need frontrunning protection,
// for example simple ERC20 transfers.
func (r *RPCRequest) doesTxNeedFrontrunningProtection(tx *ethtypes.Transaction) bool {
	gas := tx.Gas()
	r.logger.Info("[protect-check]", zap.Uint64("gas", gas))

	data := hex.EncodeToString(tx.Data())
	r.logger.Info("[protect-check] ", zap.String("tx-data", data))

	if len(data) < 8 {
		return false
	}

	if isOnFunctionWhitelist(data[0:8]) {
		return false // function being called is on our whitelist and no protection needed
	} else {
		r.logger.Info("[protect-check] Tx needs protection - function", zap.String("tx-data", data[0:8]))
		return true // needs protection if not on whitelist
	}
}
