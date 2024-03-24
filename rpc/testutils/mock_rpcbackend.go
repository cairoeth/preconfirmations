/* Package testutils Dummy RPC backend for both Ethereum node and Flashbots Relay. Implements JSON-RPC calls that the tests need. */
package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/cairoeth/preconfirmations/rpc/types"
)

var (
	MockBackendLastRawRequest              *http.Request
	MockBackendLastJSONRPCRequest          *types.JSONRPCRequest
	MockBackendLastJSONRPCRequestTimestamp time.Time
)

func MockRPCBackendReset() {
	MockBackendLastRawRequest = nil
	MockBackendLastJSONRPCRequest = nil
	MockBackendLastJSONRPCRequestTimestamp = time.Time{}
}

func handleRPCRequest(req *types.JSONRPCRequest) (result interface{}, err error) {
	MockBackendLastJSONRPCRequest = req

	switch req.Method {
	case "eth_getTransactionCount":
		if req.Params[0] == TestTxBundleFailedTooManyTimesFrom {
			return TestTxBundleFailedTooManyTimesNonce, nil
		} else if req.Params[0] == TestTxCancelAtRelayCancelFrom {
			return TestTxCancelAtRelayCancelNonce, nil
		}
		return "0x22", nil

	case "eth_call":
		return "0x12345", nil

	case "eth_getTransactionReceipt":
		if req.Params[0] == TestTxBundleFailedTooManyTimesHash {
			return nil, nil
		} else if req.Params[0] == TestTxMM2Hash {
			return nil, nil
		}

	case "eth_sendRawTransaction":
		txHash := req.Params[0].(string)
		if txHash == TestTxCancelAtRelayCancelRawTx {
			return TestTxCancelAtRelayCancelHash, nil
		}
		return "tx-hash1", nil

	case "net_version":
		return "3", nil

	case "null":
		return nil, nil

		// Relay calls
	case "eth_sendPrivateTransaction":
		param := req.Params[0].(map[string]interface{})
		if param["tx"] == TestTxBundleFailedTooManyTimesRawTx {
			return TestTxBundleFailedTooManyTimesHash, nil
		} else {
			return "tx-hash2", nil
		}

	case "eth_cancelPrivateTransaction":
		param := req.Params[0].(map[string]interface{})
		if param["txHash"] == TestTxCancelAtRelayCancelHash {
			return true, nil
		} else {
			return false, nil
		}
	}

	return "", fmt.Errorf("no RPC method handler implemented for %s", req.Method)
}

func RPCBackendHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	MockBackendLastRawRequest = req
	MockBackendLastJSONRPCRequestTimestamp = time.Now()

	log.Printf("%s %s %s\n", req.RemoteAddr, req.Method, req.URL)

	w.Header().Set("Content-Type", "application/json")
	testHeader := req.Header.Get("Test")
	w.Header().Set("Test", testHeader)

	returnError := func(id interface{}, msg string) {
		log.Println("returnError:", msg)
		res := types.JSONRPCResponse{
			ID: id,
			Error: &types.JSONRPCError{
				Code:    -32603,
				Message: msg,
			},
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("error writing response 1: %v - data: %s", err, res)
		}
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		returnError(-1, fmt.Sprintf("failed to read request body: %v", err))
		return
	}

	// Parse JSON RPC
	jsonReq := new(types.JSONRPCRequest)
	if err = json.Unmarshal(body, &jsonReq); err != nil {
		returnError(-1, fmt.Sprintf("failed to parse JSON RPC request: %v", err))
		return
	}

	rawRes, err := handleRPCRequest(jsonReq)
	if err != nil {
		returnError(jsonReq.ID, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	resBytes, err := json.Marshal(rawRes)
	if err != nil {
		fmt.Println("error mashalling rawRes:", rawRes, err)
	}

	res := types.NewJSONRPCResponse(jsonReq.ID, resBytes)

	// Write to client request
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("error writing response 2: %v - data: %s", err, rawRes)
	}
}
