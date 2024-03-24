package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/pkg/errors"

	"github.com/cairoeth/preconfirmations/rpc/types"
)

var RPCEndpointURL string // set by tests

func SendRPCAndParseResponse(req *types.JSONRPCRequest) (*types.JSONRPCResponse, error) {
	return SendRPCAndParseResponseTo(RPCEndpointURL, req)
}

func SendBatchRPCAndParseResponse(req []*types.JSONRPCRequest) ([]*types.JSONRPCResponse, error) {
	return SendBatchRPCAndParseResponseTo(RPCEndpointURL, req)
}

func SendRPCWithFastPreferenceAndParseResponse(t *testing.T, req *types.JSONRPCRequest) *types.JSONRPCResponse {
	url := RPCEndpointURL + "/fast/"
	res, err := SendRPCAndParseResponseTo(url, req)
	if err != nil {
		t.Fatal("sendRpcAndParseResponse error:", err)
	}
	return res
}

func SendRPCWithAuctionPreferenceAndParseResponse(t *testing.T, req *types.JSONRPCRequest, urlSuffix string) *types.JSONRPCResponse {
	url := RPCEndpointURL + urlSuffix
	res, err := SendRPCAndParseResponseTo(url, req)
	if err != nil {
		t.Fatal("sendRpcAndParseResponse error:", err)
	}
	return res
}

func SendRPCAndParseResponseOrFailNow(t *testing.T, req *types.JSONRPCRequest) *types.JSONRPCResponse {
	res, err := SendRPCAndParseResponse(req)
	if err != nil {
		t.Fatal("sendRpcAndParseResponse error:", err)
	}
	return res
}

func SendRPCAndParseResponseOrFailNowString(t *testing.T, req *types.JSONRPCRequest) string {
	var rpcResult string
	resp := SendRPCAndParseResponseOrFailNow(t, req)
	json.Unmarshal(resp.Result, &rpcResult)
	return rpcResult
}

func SendRPCAndParseResponseOrFailNowAllowRPCError(t *testing.T, req *types.JSONRPCRequest) *types.JSONRPCResponse {
	res, err := SendRPCAndParseResponse(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func SendRPCAndParseResponseTo(url string, req *types.JSONRPCRequest) (*types.JSONRPCResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	// fmt.Printf("%s\n", jsonData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read")
	}

	jsonRPCResp := new(types.JSONRPCResponse)

	// Check if returned an error, if so then convert to standard JSON-RPC error
	errorResp := new(types.RelayErrorResponse)
	if err := json.Unmarshal(respData, errorResp); err == nil && errorResp.Error != "" {
		// relay returned an error, convert to standard JSON-RPC error now
		jsonRPCResp.Error = &types.JSONRPCError{Message: errorResp.Error}
		return jsonRPCResp, nil
	}

	// Unmarshall JSON-RPC response and check for error inside
	if err := json.Unmarshal(respData, jsonRPCResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return jsonRPCResp, nil
}

func SendBatchRPCAndParseResponseTo(url string, req []*types.JSONRPCRequest) ([]*types.JSONRPCResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, errors.Wrap(err, "marshal")
	}

	// fmt.Printf("%s\n", jsonData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, errors.Wrap(err, "post")
	}

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read")
	}

	var jsonRPCResp []*types.JSONRPCResponse

	// Unmarshall JSON-RPC response and check for error inside
	if err := json.Unmarshal(respData, &jsonRPCResp); err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	return jsonRPCResp, nil
}
