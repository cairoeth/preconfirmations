/*
 * RPC endpoint E2E tests.
 */
package tests

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cairoeth/preconfirmations/rpc/database"

	"github.com/alicebob/miniredis"
	"github.com/cairoeth/preconfirmations/rpc/server"
	"github.com/cairoeth/preconfirmations/rpc/testutils"
	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var RPCBackendServerURL string

var relaySigningKey *ecdsa.PrivateKey

func init() {
	// var err error
	relaySigningKey, _ = crypto.HexToECDSA("7bdeed70a07d5a45546e83a88dd430f71348592e747d2d3eb23f32db003eb0e1")
	// if err != nil {
	// logger.Crit("failed to create signing key", "err", err)
	// log.Crit("failed to create signing key", "err", err)
	// }
}

// func setServerTimeNowOffset(td time.Duration) {
// 	server.Now = func() time.Time {
// 		return time.Now().Add(td)
// 	}
// }

var bundleJSONAPI *httptest.Server

// Setup RPC endpoint and mock backend servers
func testServerSetupWithMockStore() {
	db := database.NewMockStore()
	testServerSetup(db)
}

func testServerSetup(db database.Store) {
	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	// Create a fresh mock backend server (covers for both eth node and relay)
	rpcBackendServer := httptest.NewServer(http.HandlerFunc(testutils.RPCBackendHandler))
	RPCBackendServerURL = rpcBackendServer.URL
	testutils.MockRPCBackendReset()
	testutils.MockTxAPIReset()

	txAPIServer := httptest.NewServer(http.HandlerFunc(testutils.MockTxAPIHandler))
	server.ProtectTxAPIHost = txAPIServer.URL

	logger, _ := zap.NewDevelopment()

	// Create a fresh RPC endpoint server
	rpcServer, err := server.NewRPCEndPointServer(server.Configuration{
		DB:                  db,
		Logger:              logger,
		ProxyTimeoutSeconds: 10,
		ProxyURL:            RPCBackendServerURL,
		RedisURL:            redisServer.Addr(),
		RelaySigningKey:     relaySigningKey,
		RelayURL:            RPCBackendServerURL,
		Version:             "test",
	})
	if err != nil {
		panic(err)
	}
	rpcEndpointServer := httptest.NewServer(http.HandlerFunc(rpcServer.HandleHTTPRequest))
	bundleJSONAPI = httptest.NewServer(http.HandlerFunc(rpcServer.HandleBundleRequest))
	testutils.RPCEndpointURL = rpcEndpointServer.URL
}

/*
 * HTTP TESTS
 */
// Check headers: status and content-type
func TestStandardHeaders(t *testing.T) {
	testServerSetupWithMockStore()

	rpcRequest := types.NewJSONRPCRequest(1, "null", nil)
	jsonData, err := json.Marshal(rpcRequest)
	require.Nil(t, err, err)

	resp, err := http.Post(RPCBackendServerURL, "application/json", bytes.NewBuffer(jsonData))
	require.Nil(t, err, err)

	// Test for http status-code 200
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test for content-type: application/json
	contentTypeHeader := resp.Header.Get("content-type")
	assert.Equal(t, "application/json", strings.ToLower(contentTypeHeader))
}

// Check json-rpc id and version
func TestJsonRpc(t *testing.T) {
	testServerSetupWithMockStore()

	_id1 := float64(84363)
	rpcRequest := types.NewJSONRPCRequest(_id1, "null", nil)
	rpcResult := testutils.SendRPCAndParseResponseOrFailNow(t, rpcRequest)
	assert.Equal(t, _id1, rpcResult.ID)

	_id2 := "84363"
	rpcRequest2 := types.NewJSONRPCRequest(_id2, "null", nil)
	rpcResult2 := testutils.SendRPCAndParseResponseOrFailNow(t, rpcRequest2)
	assert.Equal(t, _id2, rpcResult2.ID)
	assert.Equal(t, "2.0", rpcResult2.Version)
}

/*
 * REQUEST TESTS
 */

// Test intercepting eth_call for Flashbots RPC contract
func TestEthCallIntercept(t *testing.T) {
	testServerSetupWithMockStore()
	var rpcResult string

	// eth_call intercept
	req := types.NewJSONRPCRequest(1, "eth_call", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1a",
	}})
	rpcResult = testutils.SendRPCAndParseResponseOrFailNowString(t, req)
	require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", rpcResult, "FlashRPC contract - eth_call intercept")

	// eth_call passthrough
	req2 := types.NewJSONRPCRequest(1, "eth_call", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1b",
	}})
	rpcResult = testutils.SendRPCAndParseResponseOrFailNowString(t, req2)
	require.Equal(t, "0x12345", rpcResult, "FlashRPC contract - eth_call passthrough")
}

func TestNetVersionIntercept(t *testing.T) {
	testServerSetupWithMockStore()
	var rpcResult string

	// eth_call intercept
	req := types.NewJSONRPCRequest(1, "net_version", nil)
	res, err := testutils.SendRPCAndParseResponseTo(RPCBackendServerURL, req)
	require.Nil(t, err, err)
	json.Unmarshal(res.Result, &rpcResult)
	require.Equal(t, "3", rpcResult, "net_version from backend")

	rpcResult = testutils.SendRPCAndParseResponseOrFailNowString(t, req)
	require.Nil(t, res.Error)
	require.Equal(t, "3", rpcResult, "net_version intercept")
}

// Ensure bundle response is the tx hash, not the bundle id
func TestSendBundleResponse(t *testing.T) {
	testServerSetupWithMockStore()

	// should be tx hash
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxBundleFailedTooManyTimesRawTx})
	rpcResult := testutils.SendRPCAndParseResponseOrFailNowString(t, reqSendRawTransaction)
	require.Equal(t, testutils.TestTxBundleFailedTooManyTimesHash, rpcResult)
}

func TestNull(t *testing.T) {
	testServerSetupWithMockStore()
	expectedResultRaw := `{"id":1,"result":null,"jsonrpc":"2.0"}` + "\n"

	// Build and do RPC call: "null"
	rpcRequest := types.NewJSONRPCRequest(1, "null", nil)
	jsonData, err := json.Marshal(rpcRequest)
	require.Nil(t, err, err)
	resp, err := http.Post(RPCBackendServerURL, "application/json", bytes.NewBuffer(jsonData))
	require.Nil(t, err, err)
	respData, err := io.ReadAll(resp.Body)
	require.Nil(t, err, err)

	// Check that raw result is expected
	require.Equal(t, expectedResultRaw, string(respData))

	// Parsing null results in "null":
	var jsonRPCResp types.JSONRPCResponse
	err = json.Unmarshal(respData, &jsonRPCResp)
	require.Nil(t, err, err)

	require.Equal(t, 4, len(jsonRPCResp.Result))
	require.Equal(t, json.RawMessage{110, 117, 108, 108}, jsonRPCResp.Result) // the bytes for null

	// Double-check that plain bytes are 'null'
	resultStr := string(jsonRPCResp.Result)
	require.Equal(t, "null", resultStr)
}

func TestGetTxReceiptNull(t *testing.T) {
	testServerSetupWithMockStore()

	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getTransactionReceipt", []interface{}{testutils.TestTxBundleFailedTooManyTimesHash})
	jsonResp := testutils.SendRPCAndParseResponseOrFailNow(t, reqGetTransactionCount)
	// fmt.Println(jsonResp)
	require.Equal(t, "null", string(jsonResp.Result))

	jsonResp, err := testutils.SendRPCAndParseResponseTo(RPCBackendServerURL, reqGetTransactionCount)
	require.Nil(t, err, err)

	fmt.Println(jsonResp)
	require.Equal(t, "null", string(jsonResp.Result))
}

func TestMetamaskFix(t *testing.T) {
	testServerSetupWithMockStore()
	testutils.MockTxAPIStatusForHash[testutils.TestTxMM2Hash] = types.TxStatusFailed

	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getTransactionCount", []interface{}{testutils.TestTxMM2From, "latest"})
	txCountBefore := testutils.SendRPCAndParseResponseOrFailNowString(t, reqGetTransactionCount)

	// first sendRawTransaction call: rawTx that triggers the error (creates MM cache entry)
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxMM2RawTx})
	r1 := testutils.SendRPCAndParseResponseOrFailNowAllowRPCError(t, reqSendRawTransaction)
	require.Nil(t, r1.Error, r1.Error)
	fmt.Printf("\n\n\n\n\n")

	// call getTxReceipt to trigger query to Tx API
	reqGetTransactionReceipt := types.NewJSONRPCRequest(1, "eth_getTransactionReceipt", []interface{}{testutils.TestTxMM2Hash})
	jsonResp := testutils.SendRPCAndParseResponseOrFailNow(t, reqGetTransactionReceipt)
	require.Nil(t, jsonResp.Error)
	require.Equal(t, "null", string(jsonResp.Result))
	// require.Equal(t, "Transaction failed", jsonResp.Error.Message)

	// At this point, the tx hash should be blacklisted and too high a nonce is returned
	valueAfter1 := testutils.SendRPCAndParseResponseOrFailNowString(t, reqGetTransactionCount)
	require.NotEqual(t, txCountBefore, valueAfter1)
	require.Equal(t, "0x3b9aca01", valueAfter1)

	// getTransactionCount 2/4 should return the same (fixed) value
	valueAfter2 := testutils.SendRPCAndParseResponseOrFailNowString(t, reqGetTransactionCount)
	require.Equal(t, valueAfter1, valueAfter2)

	// getTransactionCount 3/4 should return the same (fixed) value
	valueAfter3 := testutils.SendRPCAndParseResponseOrFailNowString(t, reqGetTransactionCount)
	require.Equal(t, valueAfter1, valueAfter3)

	// getTransactionCount 4/4 should return the same (fixed) value
	valueAfter4 := testutils.SendRPCAndParseResponseOrFailNowString(t, reqGetTransactionCount)
	require.Equal(t, valueAfter1, valueAfter4)

	// getTransactionCount 5 should return the initial value
	valueAfter5 := testutils.SendRPCAndParseResponseOrFailNowString(t, reqGetTransactionCount)
	require.Equal(t, txCountBefore, valueAfter5)
}

func TestRelayTx(t *testing.T) {
	testServerSetupWithMockStore()

	// sendRawTransaction adds tx to MM cache entry, to be used at later eth_getTransactionReceipt call
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxBundleFailedTooManyTimesRawTx})
	r1 := testutils.SendRPCAndParseResponseOrFailNowAllowRPCError(t, reqSendRawTransaction)
	require.Nil(t, r1.Error)

	// Ensure that request called eth_sendPrivateTransaction with correct param
	require.Equal(t, "eth_sendPrivateTransaction", testutils.MockBackendLastJSONRPCRequest.Method)

	resp := testutils.MockBackendLastJSONRPCRequest.Params[0].(map[string]interface{})
	require.Equal(t, testutils.TestTxBundleFailedTooManyTimesRawTx, resp["tx"])

	// Ensure that request was signed properly
	pubkey := crypto.PubkeyToAddress(relaySigningKey.PublicKey).Hex()
	require.Equal(t, pubkey+":0xced5faaf1709075e73c8199aa32fe48f46d7fef3fe2bffd8f6bfb2953add6e89373977dd4ab50e494548944519b8cef1c4dfa810c60fda7025b50e4c324e98d100", testutils.MockBackendLastRawRequest.Header.Get("X-Flashbots-Signature"))

	// Check result - should be the tx hash
	var res string
	json.Unmarshal(r1.Result, &res)
	require.Equal(t, testutils.TestTxBundleFailedTooManyTimesHash, res)

	timeStampFirstRequest := testutils.MockBackendLastJSONRPCRequestTimestamp

	// Send tx again, should not arrive at backend
	testutils.SendRPCAndParseResponseOrFailNowAllowRPCError(t, reqSendRawTransaction)
	require.Nil(t, r1.Error)
	require.Equal(t, timeStampFirstRequest, testutils.MockBackendLastJSONRPCRequestTimestamp)

	// Ensure nonce is saved to redis
	nonce, found, err := server.RState.GetSenderMaxNonce(testutils.TestTxBundleFailedTooManyTimesFrom)
	require.Nil(t, err, err)
	require.True(t, found)
	require.Equal(t, uint64(30), nonce)
}

func TestRelayTxWithAuctionPreference(t *testing.T) {
	// Store setup
	memStore := database.NewMemStore()

	// Server setup
	testServerSetup(memStore)

	tx := testutils.TestTxBundleFailedTooManyTimesRawTx
	// sendRawTransaction adds tx to MM cache entry, to be used at later eth_getTransactionReceipt call
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{tx})
	// call rpc with auction preference
	r1 := testutils.SendRPCWithAuctionPreferenceAndParseResponse(t, reqSendRawTransaction, "/?hint=calldata&hint=contract_address")
	require.Nil(t, r1.Error)

	// Ensure that request called eth_sendPrivateTransaction with correct param
	require.Equal(t, "eth_sendPrivateTransaction", testutils.MockBackendLastJSONRPCRequest.Method)

	resp := testutils.MockBackendLastJSONRPCRequest.Params[0].(map[string]interface{})
	require.Equal(t, tx, resp["tx"])
	// Ensure fast endpoint is called and fast preference is set
	auctionPref := resp["preferences"].(map[string]interface{})["privacy"].(map[string]interface{})
	require.NotNil(t, auctionPref)

	expectedHints := []string{"calldata", "contract_address", "special_logs"}
	hintPref := auctionPref["hints"].([]interface{})
	for i, hint := range hintPref {
		strHint := hint.(string)
		require.Equal(t, expectedHints[i], strHint)
	}

	require.Equal(t, 1, len(memStore.EthSendRawTxs))
}

func TestRelayTxWithIncorrectAuctionPreference(t *testing.T) {
	// Store setup
	memStore := database.NewMemStore()

	// Server setup
	testServerSetup(memStore)

	tx := testutils.TestTxBundleFailedTooManyTimesRawTx
	// sendRawTransaction adds tx to MM cache entry, to be used at later eth_getTransactionReceipt call
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{tx})
	// call rpc with auction preference
	r1 := testutils.SendRPCWithAuctionPreferenceAndParseResponse(t, reqSendRawTransaction, "/?hint=incorrect")
	require.Contains(t, r1.Error.Message, "Incorrect auction hint")
}

func TestRelayCancelTx(t *testing.T) {
	testServerSetupWithMockStore()

	// sendRawTransaction of the initial TX
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxCancelAtRelayInitialRawTx})
	testutils.SendRPCAndParseResponseOrFailNow(t, reqSendRawTransaction)

	// Ensure that request called eth_sendPrivateTransaction on the Relay
	require.Equal(t, "eth_sendPrivateTransaction", testutils.MockBackendLastJSONRPCRequest.Method)

	// Ensure that the RPC backend sent the rawTx to the relay
	resp := testutils.MockBackendLastJSONRPCRequest.Params[0].(map[string]interface{})
	require.Equal(t, testutils.TestTxCancelAtRelayInitialRawTx, resp["tx"])

	// Send cancel-tx to the RPC backend
	reqCancelTx := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxCancelAtRelayCancelRawTx})
	cancelResp := testutils.SendRPCAndParseResponseOrFailNow(t, reqCancelTx)

	// Ensure that request called eth_sendPrivateTransaction on the Relay
	require.Equal(t, "eth_cancelPrivateTransaction", testutils.MockBackendLastJSONRPCRequest.Method)
	var res string
	json.Unmarshal(cancelResp.Result, &res)

	// Ensure the response is the tx hash
	require.Equal(t, testutils.TestTxCancelAtRelayCancelHash, res)
}

// cancel-tx without initial related tx would just go to mempool
func TestRelayCancelTxWithoutInitialTx(t *testing.T) {
	testServerSetupWithMockStore()

	// Send cancel-tx to the RPC backend
	reqCancelTx := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxCancelAtRelayCancelRawTx})
	cancelResp := testutils.SendRPCAndParseResponseOrFailNow(t, reqCancelTx)

	// Ensure that request called eth_sendRawTransaction on the mempool node, instead of eth_sendPrivateTransaction on the Relay
	// (since no valid initial tx was found)
	require.Equal(t, "eth_sendRawTransaction", testutils.MockBackendLastJSONRPCRequest.Method)
	var res string
	json.Unmarshal(cancelResp.Result, &res)

	// Ensure the response is the tx hash
	require.Equal(t, testutils.TestTxCancelAtRelayCancelHash, res)
}

// tx with wrong nonce should be rejected
func TestRelayTxWithWrongNonce(t *testing.T) {
	testServerSetupWithMockStore()

	nonceOrig := testutils.TestTxBundleFailedTooManyTimesNonce
	testutils.TestTxBundleFailedTooManyTimesNonce = "0x1f"
	defer func() { testutils.TestTxBundleFailedTooManyTimesNonce = nonceOrig }()

	// Send cancel-tx to the RPC backend
	req1 := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxBundleFailedTooManyTimesRawTx})
	resp1 := testutils.SendRPCAndParseResponseOrFailNow(t, req1)

	// Ensure the response has an error
	require.NotNil(t, resp1.Error)
	require.Equal(t, "invalid nonce", resp1.Error.Message)
}

// Test batch request with multiple eth raw transaction
func TestBatch_eth_sendRawTransaction(t *testing.T) {
	t.Skip()
	testServerSetupWithMockStore()

	var batch []*types.JSONRPCRequest
	for i := range "testing" {
		rpcRequest := types.NewJSONRPCRequest(i, "eth_sendRawTransaction", []interface{}{testutils.TestTxCancelAtRelayCancelRawTx})
		batch = append(batch, rpcRequest)
	}
	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 7)
}

// Test batch request with different eth transaction
func TestBatch_eth_transaction(t *testing.T) {
	t.Skip()
	testServerSetupWithMockStore()

	var batch []*types.JSONRPCRequest
	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getTransactionCount", []interface{}{testutils.TestTxMM2From, "latest"})
	batch = append(batch, reqGetTransactionCount)
	// first sendRawTransaction call: rawTx that triggers the error (creates MM cache entry)
	reqSendRawTransaction := types.NewJSONRPCRequest(2, "eth_sendRawTransaction", []interface{}{testutils.TestTxMM2RawTx})
	batch = append(batch, reqSendRawTransaction)
	// call getTxReceipt to trigger query to Tx API
	reqGetTransactionReceipt := types.NewJSONRPCRequest(3, "eth_getTransactionReceipt", []interface{}{testutils.TestTxMM2Hash})
	batch = append(batch, reqGetTransactionReceipt)

	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 3)

	m := map[float64]*types.JSONRPCResponse{
		float64(1): {ID: float64(1), Result: []byte(`"0x22"`), Error: nil, Version: "2.0"},
		float64(2): {ID: float64(2), Result: []byte(`"tx-hash1"`), Error: nil, Version: "2.0"},
		float64(3): {ID: float64(3), Result: []byte(`null`), Error: nil, Version: "2.0"},
	}
	for _, j := range res {
		assert.Equal(t, m[j.ID.(float64)], j)
	}
}

// Test batch request with different eth transaction
func TestBatch_eth_call(t *testing.T) {
	t.Skip()
	testServerSetupWithMockStore()

	var batch []*types.JSONRPCRequest
	// eth_call intercept
	req := types.NewJSONRPCRequest(1, "eth_call", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1a",
	}})
	batch = append(batch, req)
	// eth_call passthrough
	req2 := types.NewJSONRPCRequest(2, "eth_call", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1b",
	}})
	batch = append(batch, req2)
	reqGetTransactionCount := types.NewJSONRPCRequest(3, "eth_getTransactionCount", []interface{}{testutils.TestTxMM2From, "latest"})
	batch = append(batch, reqGetTransactionCount)
	// first sendRawTransaction call: rawTx that triggers the error (creates MM cache entry)
	reqSendRawTransaction := types.NewJSONRPCRequest(4, "eth_sendRawTransaction", []interface{}{testutils.TestTxMM2RawTx})
	batch = append(batch, reqSendRawTransaction)
	// call getTxReceipt to trigger query to Tx API
	reqGetTransactionReceipt := types.NewJSONRPCRequest(5, "eth_getTransactionReceipt", []interface{}{testutils.TestTxMM2Hash})
	batch = append(batch, reqGetTransactionReceipt)

	m := map[float64]*types.JSONRPCResponse{
		float64(1): {ID: float64(1), Result: []byte(`"0x0000000000000000000000000000000000000000000000000000000000000001"`), Error: nil, Version: "2.0"},
		float64(2): {ID: float64(2), Result: []byte(`"0x12345"`), Error: nil, Version: "2.0"},
		float64(3): {ID: float64(3), Result: []byte(`"0x22"`), Error: nil, Version: "2.0"},
		float64(4): {ID: float64(4), Result: []byte(`"tx-hash1"`), Error: nil, Version: "2.0"},
		float64(5): {ID: float64(5), Result: []byte(`null`), Error: nil, Version: "2.0"},
	}
	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 5)
	for _, j := range res {
		assert.Equal(t, m[j.ID.(float64)], j)
	}
}

// Test batch request with different transaction
func TestBatch_CombinationOfSuccessAndFailure(t *testing.T) {
	t.Skip()
	testServerSetupWithMockStore()

	var batch []*types.JSONRPCRequest
	// eth_call intercept
	req := types.NewJSONRPCRequest(1, "eth_call", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1a",
	}})
	batch = append(batch, req)
	// eth_call passthrough
	req2 := types.NewJSONRPCRequest(1, "eth_callxvssfa", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1b",
	}})
	batch = append(batch, req2)
	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getTransactionCount", []interface{}{testutils.TestTxMM2From, "latest"})
	batch = append(batch, reqGetTransactionCount)
	// first sendRawTransaction call: rawTx that triggers the error (creates MM cache entry)
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransactionxxx", []interface{}{testutils.TestTxMM2RawTx})
	batch = append(batch, reqSendRawTransaction)
	// call getTxReceipt to trigger query to Tx API
	reqGetTransactionReceipt := types.NewJSONRPCRequest(1, "eth_getTransactionReceipt", []interface{}{testutils.TestTxMM2Hash})
	batch = append(batch, reqGetTransactionReceipt)

	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 5)
}

// Test batch request with multiple eth raw transaction
func TestBatch_Validate_eth_sendRawTransaction_Error(t *testing.T) {
	t.Skip()
	testServerSetupWithMockStore()
	// key=request-id, value=json-rpc error
	m := map[float64]int{
		1: types.JSONRPCInvalidParams,
		2: types.JSONRPCInvalidParams,
		3: types.JSONRPCInvalidParams,
		4: types.JSONRPCInvalidRequest,
	}
	var batch []*types.JSONRPCRequest

	r1 := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{})     // no params
	r2 := types.NewJSONRPCRequest(2, "eth_sendRawTransaction", nil)                 // nil params
	r3 := types.NewJSONRPCRequest(3, "eth_sendRawTransaction", []interface{}{"x"})  // invalid params
	r4 := types.NewJSONRPCRequest(4, "eth_sendRawTransaction", []interface{}{"xy"}) // invalid request
	batch = append(batch, r1, r2, r3, r4)

	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 4)
	for _, r := range res {
		assert.Equal(t, m[r.ID.(float64)], r.Error.Code)
	}
}

// Whitehat Tests
func TestWhitehatBundleCollection(t *testing.T) {
	testServerSetupWithMockStore()

	bundleID := "123"
	url := testutils.RPCEndpointURL + "?bundle=" + bundleID

	// sendRawTransaction adds tx to MM cache entry, to be used at later eth_getTransactionReceipt call
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxBundleFailedTooManyTimesRawTx})
	resp, err := testutils.SendRPCAndParseResponseTo(url, reqSendRawTransaction)
	require.Nil(t, err, err)
	require.Nil(t, resp.Error, resp.Error)

	// Last request should be network version (executed on start)
	require.Equal(t, &types.JSONRPCRequest{ID: float64(1), Method: "net_version", Params: []interface{}{}, Version: "2.0"}, testutils.MockBackendLastJSONRPCRequest)
	// Check redis
	txs, err := server.RState.GetWhitehatBundleTx(bundleID)
	require.Nil(t, err, err)
	require.Equal(t, 1, len(txs))

	// Send again (#2)
	resp, err = testutils.SendRPCAndParseResponseTo(url, reqSendRawTransaction)
	require.Nil(t, err, err)
	require.Nil(t, resp.Error, resp.Error)

	// Check redis (#2)
	txs, err = server.RState.GetWhitehatBundleTx(bundleID)
	require.Nil(t, err, err)
	require.Equal(t, 1, len(txs))

	// Check JSON API
	jsonAPIURL := bundleJSONAPI.URL + "/bundle?id=" + bundleID
	fmt.Println("jsonAPIURL: ", jsonAPIURL)
	res, err := http.Get(jsonAPIURL)
	require.Nil(t, err, err)
	body, err := io.ReadAll(res.Body)
	require.Nil(t, err, err)
	fmt.Println(string(body))
	bundleResponse := new(types.BundleResponse)
	err = json.Unmarshal(body, bundleResponse)
	require.Nil(t, err, err)
	require.Equal(t, bundleID, bundleResponse.BundleID)
	require.Equal(t, 1, len(bundleResponse.RawTxs))
}

func TestWhitehatBundleCollectionGetBalance(t *testing.T) {
	testServerSetupWithMockStore()
	bundleID := "123"
	url := testutils.RPCEndpointURL + "?bundle=" + bundleID

	// sendRawTransaction adds tx to MM cache entry, to be used at later eth_getTransactionReceipt call
	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getBalance", []interface{}{testutils.TestTxMM2From, "latest"})
	resp, err := testutils.SendRPCAndParseResponseTo(url, reqGetTransactionCount)
	require.Nil(t, err, err)
	require.Nil(t, resp.Error, resp.Error)
	val := ""
	err = json.Unmarshal(resp.Result, &val)
	require.Nil(t, err, err)
	require.Equal(t, "0x56bc75e2d63100000", val)
}

func Test_StoreRequests(t *testing.T) {
	// Store setup
	memStore := database.NewMemStore()

	// Server setup
	testServerSetup(memStore)

	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getTransactionReceipt", []interface{}{testutils.TestTxBundleFailedTooManyTimesHash})
	_ = testutils.SendRPCAndParseResponseOrFailNow(t, reqGetTransactionCount)
	// sendRawTransaction of the initial TX
	reqSendRawTransaction1 := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxCancelAtRelayInitialRawTx})
	testutils.SendRPCAndParseResponseOrFailNow(t, reqSendRawTransaction1)

	// sendRawTransaction adds tx to MM cache entry, to be used at later eth_getTransactionReceipt call
	reqSendRawTransaction2 := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxBundleFailedTooManyTimesRawTx})
	r1 := testutils.SendRPCAndParseResponseOrFailNowAllowRPCError(t, reqSendRawTransaction2)
	require.Nil(t, r1.Error)

	require.Equal(t, 2, len(memStore.Requests))
	require.Equal(t, 2, len(memStore.EthSendRawTxs))
	for _, txs := range memStore.EthSendRawTxs {
		for _, tx := range txs {
			assert.Equal(t, true, tx.NeedsFrontRunningProtection)
		}
	}
}

func Test_StoreBatchRequests(t *testing.T) {
	t.Skip()
	// Store setup
	memStore := database.NewMemStore()
	// Server setup
	testServerSetup(memStore)

	var batch []*types.JSONRPCRequest
	// eth_call intercept
	req := types.NewJSONRPCRequest(1, "eth_call", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1a",
	}})
	batch = append(batch, req)
	// eth_call passthrough
	req2 := types.NewJSONRPCRequest(1, "eth_callxvssfa", []interface{}{map[string]string{
		"from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		"to":   "0xf1a54b0759b58661cea17cff19dd37940a9b5f1b",
	}})
	batch = append(batch, req2)
	reqGetTransactionCount := types.NewJSONRPCRequest(1, "eth_getTransactionCount", []interface{}{testutils.TestTxMM2From, "latest"})
	batch = append(batch, reqGetTransactionCount)
	// first sendRawTransaction call: rawTx that triggers the error (creates MM cache entry)
	reqSendRawTransaction := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxMM2RawTx})
	batch = append(batch, reqSendRawTransaction)
	// call getTxReceipt to trigger query to Tx API
	reqGetTransactionReceipt := types.NewJSONRPCRequest(1, "eth_getTransactionReceipt", []interface{}{testutils.TestTxMM2Hash})
	batch = append(batch, reqGetTransactionReceipt)

	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 5)
	require.Equal(t, 1, len(memStore.Requests))
	require.Equal(t, 1, len(memStore.EthSendRawTxs))
}

func Test_StoreValidateTxs(t *testing.T) {
	t.Skip()
	// Store setup
	memStore := database.NewMemStore()

	// Server setup
	testServerSetup(memStore)

	var batch []*types.JSONRPCRequest

	// call sendRawTx
	reqSendRawTransactionInvalidNonce1 := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxInvalidNonce1})
	batch = append(batch, reqSendRawTransactionInvalidNonce1)

	reqSendRawTransactionInvalidNonce2 := types.NewJSONRPCRequest(1, "eth_sendRawTransaction", []interface{}{testutils.TestTxInvalidNonce2})
	batch = append(batch, reqSendRawTransactionInvalidNonce2)

	res, err := testutils.SendBatchRPCAndParseResponse(batch)
	require.Nil(t, err, err)
	assert.Equal(t, len(res), 2)
	require.Equal(t, 1, len(memStore.Requests))
	require.Equal(t, 1, len(memStore.EthSendRawTxs))

	for _, entries := range memStore.EthSendRawTxs {
		for _, entry := range entries {
			require.True(t, entry.NeedsFrontRunningProtection)
			require.Equal(t, "invalid nonce", entry.Error)
			require.Equal(t, -32603, entry.ErrorCode)
			require.Equal(t, 10, len(entry.TxSmartContractMethod))
			require.False(t, entry.Fast)
		}
	}
}
