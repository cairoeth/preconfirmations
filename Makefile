VERSION := $(shell git describe --tags --always --dirty="-dev")

all: clean build

v:
	@echo "Version: ${VERSION}"

clean:
	rm -rf build/

# Precon-Share
build-share:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/node precon-share/cmd/node/main.go

# Precon-Operator
build-operator:
	go build -trimpath -ldflags "-X main.version=${VERSION}" -v -o build/operator precon-operator/cmd/operator/main.go

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
