BINARY := glox

.PHONY: all build run clean fmt vet tidy test

all: build

build:
	go build -o bin/$(BINARY) ./cmd/glox

run: build
	./bin/$(BINARY)

clean:
	rm -rf bin

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

test:
	go test ./...
