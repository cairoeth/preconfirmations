package server

import (
	"encoding/json"
	"net/http"

	"github.com/cairoeth/preconfirmations/rpc/types"
	"go.uber.org/zap"
)

func (r *RPCRequestHandler) writeHeaderContentTypeJSON() {
	(*r.respw).Header().Set("Content-Type", "application/json")
}

func (r *RPCRequestHandler) _writeRPCResponse(res *types.JSONRPCResponse) {
	// If the request is single and not batch
	// Write content type
	r.writeHeaderContentTypeJSON() // Set content type to json
	(*r.respw).WriteHeader(http.StatusOK)
	// Write response
	if err := json.NewEncoder(*r.respw).Encode(res); err != nil {
		r.logger.Error("[_writeRPCResponse] Failed writing rpc response", zap.Error(err))
		(*r.respw).WriteHeader(http.StatusInternalServerError)
	}
}

//lint:ignore U1000 Ignore all unused code, it's generated
func (r *RPCRequestHandler) _writeRPCBatchResponse(res []*types.JSONRPCResponse) {
	r.writeHeaderContentTypeJSON() // Set content type to json
	(*r.respw).WriteHeader(http.StatusOK)
	// Write response
	if err := json.NewEncoder(*r.respw).Encode(res); err != nil {
		r.logger.Error("[_writeRPCBatchResponse] Failed writing rpc response", zap.Error(err))
		(*r.respw).WriteHeader(http.StatusInternalServerError)
	}
}
