# 🔌 Preconfirmations AVS 


## Preconf-Share

### Dependencies

- Redis: Used for hint streaming and priority queue.
- Postgres: Used for storing bundles and historical hints.

### Configuration

The full list of configuration options can be found in [preconf-share/cmd/node/main.go](preconf-share/cmd/node/main.go).

### Running Locally

```bash
docker-compose up # start services: redis and postgres

# apply migration
for file in preconf-share/sql/*.sql; do psql "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -f $file; done

# run blockchain
anvil

# run node
make && ./build/node
```

## Preconf-Contract

### Deploy

```bash
forge script script/PreconfDeploy.s.sol:PreconfDeploy --rpc-url $RPC_URL --broadcast -vvvv
```