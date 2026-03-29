BINARY_NAME=pg-util

build:
	go build -o $(BINARY_NAME) ./cmd/main.go

clean:
	rm -f $(BINARY_NAME)

test:
	go test -v ./integration

.PHONY: build clean test