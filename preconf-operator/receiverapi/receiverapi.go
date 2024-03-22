package receiverapi

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cairoeth/preconfirmations-avs/rpc/types"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

type ReceiveApi struct {
	ipPortAddr    string
	logger        logging.Logger
}

type ReceiveResponse struct {
	jsonrpc   string `json:"jsonrpc"`
    id    	  int    `json:"id"`
}

func NewReceiveApi(IpPortAddr string, logger logging.Logger) *ReceiveApi {
	receiveApi := &ReceiveApi{
		ipPortAddr:    IpPortAddr,
		logger:        logger,
	}
	return receiveApi
}

// Start starts the receiver api server in a goroutine
func (api *ReceiveApi) Start() <-chan error {
	api.logger.Infof("Starting receiver api server at address %v", api.ipPortAddr)

	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    api.ipPortAddr,
		Handler: mux,
	}

	mux.HandleFunc("/receive", api.receive)
	errChan := run(api.logger, &httpServer)
	return errChan
}

func (api *ReceiveApi) receive(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
    var txs types.JsonRpcRequest
    err := decoder.Decode(&txs)
    if err != nil {
        api.logger.Error("could not read request body", "err", err)
    }
    api.logger.Info("processed receive", "txs", txs)

	


	response := &ReceiveResponse{"2.0", 1}
	err = jsonResponse(w, response)
	if err != nil {
		api.logger.Error("Error in receive endpoint", "err", err)
	}
}

func run(logger logging.Logger, httpServer *http.Server) <-chan error {
	errChan := make(chan error, 1)
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()

		logger.Info("shutdown signal received")

		defer func() {
			stop()
			close(errChan)
		}()

		if err := httpServer.Shutdown(context.Background()); err != nil {
			errChan <- err
		}
		logger.Info("shutdown completed")
	}()

	go func() {
		logger.Info("receiver api server running", "addr", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	return errChan
}

func jsonResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	return err
}
