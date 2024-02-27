## Preconf RPC Endpoint

## Usage

Wallet RPC for preconfirmations, wrapping on top of a network.

To run the server, run the following command:

```bash
go run cmd/server/main.go -redis REDIS_URL -signingKey ETH_PRIVATE_KEY -proxy PROXY_URL

# For development, you can use built-in redis and create a random signing key
go run cmd/server/main.go -redis dev -signingKey dev -proxy PROXY_URL

# You can use the DEBUG_DONT_SEND_RAWTX to skip sending transactions anywhere (useful for local testing):
DEBUG_DONT_SEND_RAWTX=1 go run cmd/server/main.go -redis dev -signingKey dev -proxy PROXY_URL
```

Example Single request:

```bash
curl localhost:9000 -f -d '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["latest", false],"id":1}'
```
