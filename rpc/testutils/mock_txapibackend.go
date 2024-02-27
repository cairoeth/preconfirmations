package testutils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/cairoeth/preconfirmations-avs/rpc/types"
)

var MockTxApiStatusForHash map[string]types.PrivateTxStatus = make(map[string]types.PrivateTxStatus)

func MockTxApiReset() {
	MockTxApiStatusForHash = make(map[string]types.PrivateTxStatus)
}

func MockTxApiHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	fmt.Println("TX API", req.URL)

	if !strings.HasPrefix(req.URL.Path, "/tx/") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	txHash := req.URL.Path[4:] // by default, the first 4 characters are "/tx/"
	resp := types.PrivateTxApiResponse{Status: types.TxStatusUnknown}

	if status, found := MockTxApiStatusForHash[txHash]; found {
		resp.Status = status
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error writing response 2: %v - data: %v", err, resp)
	}
}
