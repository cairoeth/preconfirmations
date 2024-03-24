package server

import (
	"crypto/ecdsa"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/cairoeth/preconfirmations/rpc/database"
	"github.com/cairoeth/preconfirmations/rpc/types"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// RPCRequestHandler RPC request handler for a single/ batch JSON-RPC request
type RPCRequestHandler struct {
	respw               *http.ResponseWriter
	req                 *http.Request
	logger              *zap.Logger
	timeStarted         time.Time
	defaultProxyURL     string
	proxyTimeoutSeconds int
	relaySigningKey     *ecdsa.PrivateKey
	relayURL            string
	uid                 uuid.UUID
	requestRecord       *requestRecord
	builderNames        []string
	chainID             []byte
}

func NewRPCRequestHandler(logger *zap.Logger, respw *http.ResponseWriter, req *http.Request, proxyURL string, proxyTimeoutSeconds int, relaySigningKey *ecdsa.PrivateKey, relayURL string, db database.Store, builderNames []string, chainID []byte) *RPCRequestHandler {
	return &RPCRequestHandler{
		logger:              logger,
		respw:               respw,
		req:                 req,
		timeStarted:         Now(),
		defaultProxyURL:     proxyURL,
		proxyTimeoutSeconds: proxyTimeoutSeconds,
		relaySigningKey:     relaySigningKey,
		relayURL:            relayURL,
		uid:                 uuid.New(),
		requestRecord:       NewRequestRecord(db),
		builderNames:        builderNames,
		chainID:             chainID,
	}
}

// nolint
func (r *RPCRequestHandler) process() {
	// r.logger = r.logger.With("uid", r.uid)
	r.logger.Info("[process] POST request received")

	defer r.finishRequest()
	r.requestRecord.requestEntry.ReceivedAt = r.timeStarted
	r.requestRecord.requestEntry.ID = r.uid
	r.requestRecord.UpdateRequestEntry(r.req, http.StatusOK, "")

	whitehatBundleID := r.req.URL.Query().Get("bundle")
	isWhitehatBundleCollection := whitehatBundleID != ""

	origin := r.req.Header.Get("Origin")
	referer := r.req.Header.Get("Referer")

	// If users specify a proxy url in their rpc endpoint they can have their requests proxied to that endpoint instead of Infura
	// e.g. https://rpc.flashbots.net?url=http://RPC-ENDPOINT.COM
	customProxyURL, ok := r.req.URL.Query()["url"]
	if ok && len(customProxyURL[0]) > 1 {
		r.defaultProxyURL = customProxyURL[0]
		r.logger.Info("[process] Using custom url", zap.String("url", r.defaultProxyURL))
	}

	// Decode request JSON RPC
	defer r.req.Body.Close()
	body, err := io.ReadAll(r.req.Body)
	if err != nil {
		r.requestRecord.UpdateRequestEntry(r.req, http.StatusBadRequest, err.Error())
		r.logger.Error("[process] Failed to read request body", zap.Error(err))
		(*r.respw).WriteHeader(http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		r.requestRecord.UpdateRequestEntry(r.req, http.StatusBadRequest, "empty request body")
		(*r.respw).WriteHeader(http.StatusBadRequest)
		return
	}

	// create rpc proxy client for making proxy request
	client := NewRPCProxyClient(r.logger, r.defaultProxyURL, r.proxyTimeoutSeconds)

	r.requestRecord.UpdateRequestEntry(r.req, http.StatusOK, "") // Data analytics

	// Parse JSON RPC payload
	var jsonReq *types.JSONRPCRequest
	if err = json.Unmarshal(body, &jsonReq); err != nil {
		r.logger.Warn("[process] Parse payload", zap.Error(err))
		(*r.respw).WriteHeader(http.StatusBadRequest)
		return
	}

	// mev-share parameters
	urlParams, err := ExtractParametersFromURL(r.req.URL, r.builderNames)
	if err != nil {
		r.logger.Warn("[process] Invalid auction preference", zap.Error(err))
		res := AuctionPreferenceErrorToJSONRPCResponse(jsonReq, err)
		r._writeRPCResponse(res)
		return
	}
	// Process single request
	r.processRequest(client, jsonReq, origin, referer, isWhitehatBundleCollection, r.defaultProxyURL, urlParams)
}

// processRequest handles single request
func (r *RPCRequestHandler) processRequest(client RPCProxyClient, jsonReq *types.JSONRPCRequest, origin, referer string, isWhitehatBundleCollection bool, whitehatBundleID string, urlParams URLParameters) {
	var entry *database.EthSendRawTxEntry
	if jsonReq.Method == "eth_sendRawTransaction" {
		entry = r.requestRecord.AddEthSendRawTxEntry(uuid.New())
	}
	// Handle single request
	rpcReq := NewRPCRequest(r.logger, client, jsonReq, r.relaySigningKey, r.relayURL, origin, referer, isWhitehatBundleCollection, whitehatBundleID, entry, urlParams, r.chainID)
	res := rpcReq.ProcessRequest()
	// Write response
	r._writeRPCResponse(res)
}

func (r *RPCRequestHandler) finishRequest() {
	reqDuration := time.Since(r.timeStarted) // At end of request, log the time it needed
	r.requestRecord.requestEntry.RequestDurationMs = reqDuration.Milliseconds()
	go func() {
		// Save both request entry and raw tx entries if present
		if err := r.requestRecord.SaveRecord(); err != nil {
			r.logger.Error("saveRecord failed", zap.Error(err))
		}
	}()
	r.logger.Info("Request finished", zap.Int("duration", int(reqDuration.Seconds())))
}
