package server

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type RPCProxyClient interface {
	ProxyRequest(body []byte) (*http.Response, error)
}

type rpcProxyClient struct {
	logger     *zap.Logger
	httpClient http.Client
	proxyURL   string
}

func NewRPCProxyClient(logger *zap.Logger, proxyURL string, timeoutSeconds int) RPCProxyClient {
	return &rpcProxyClient{
		logger:     logger,
		httpClient: http.Client{Timeout: time.Second * time.Duration(timeoutSeconds)},
		proxyURL:   proxyURL,
	}
}

// ProxyRequest using http client to make http post request
func (n *rpcProxyClient) ProxyRequest(body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, n.proxyURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	start := time.Now()
	res, err := n.httpClient.Do(req)
	n.logger.Info("[ProxyRequest] completed", zap.Duration("timeNeeded", time.Since(start)))
	return res, err
}
