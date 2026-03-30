BINARY_NAME=pg-util

build: deps
	go build -o $(BINARY_NAME) ./cmd/main.go

clean:
	rm -f $(BINARY_NAME)

deps:
	go mod tidy

test: deps
	go clean -testcache
	go test -v ./integration

.PHONY: build clean test