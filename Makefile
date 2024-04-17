VERSION := $(shell git describe --tags --always --dirty="-dev")

v:
	@echo "Version: ${VERSION}"

clean:
	rm -rf build/

# Build preconf-share
build-share:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/node preconf-share/cmd/node/main.go

# Build preconf-operator
build-operator:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/operator preconf-operator/cmd/operator/main.go

# Build rpc
build-rpc:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/rpc rpc/cmd/server/main.go

# Build all components
build-all: clean build-share build-operator build-rpc

# Run preconf-share
run-share:
	make build-share
	cd preconf-share && docker-compose rm -f -s && docker compose up -d --force-recreate --build --wait && sleep 1
	for file in preconf-share/sql/*.sql; do psql "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -f $$file; done
	./build/node

# Run preconf-operator
run-operator:
	make build-operator
	./build/operator --config preconf-operator/config.yaml

# Run anvil with prepared state
run-anvil:
	anvil --load-state contracts/anvil-state.json --silent

# Run the rpc
run-rpc:
	go run rpc/cmd/server/main.go -redis dev -signingKey dev -proxy http://127.0.0.1:8545

# Run all components concurrently
run-all:
	make -j run-anvil run-share run-operator run-rpc &

# Run all in background and example (used by CI)
run-ci:
	make run-all && sleep 200 && python test_tx.py && exit 0

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	gofmt -d -s .
	gofumpt -d -extra .
	go vet ./...
	go list ./... | grep -F -e contracts/ -v | xargs staticcheck
	# golangci-lint run

fmt:
	gofumpt -l -w -extra .
