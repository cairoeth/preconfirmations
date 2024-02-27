package main

import (
	"fmt"
	"net/http"

	"github.com/cairoeth/preconfirmations-avs/rpc/testutils"
)

func main() {
	port := 8090
	http.HandleFunc("/", testutils.RpcBackendHandler)
	fmt.Printf("rpc backend listening on localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}
