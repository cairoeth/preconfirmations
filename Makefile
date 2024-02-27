VERSION := $(shell git describe --tags --always --dirty="-dev")

all: clean build

v:
	@echo "Version: ${VERSION}"

clean:
	rm -rf build/

# Preconf-Share
build-share:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/node preconf-share/cmd/node/main.go

# Preconf-Operator
build-operator:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/operator preconf-operator/cmd/operator/main.go

# RPC
build-rpc:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/rpc rpc/cmd/server/main.go

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	gofmt -d -s .
	gofumpt -d -extra .
	go vet ./...
	staticcheck ./...
	golangci-lint run

fmt:
	gofumpt -l -w -extra .
