BIN     := sharelock
BIN_SRV := sharelock-server
CERT    := cert.pem
KEY     := key.pem

.PHONY: all build client server clean test vet bench gen-cert help

all: build gen-cert

build: client server

client:
	go build -o $(BIN) ./cmd/client

server:
	go build -o $(BIN_SRV) ./cmd/server

gen-cert:
	@if [ ! -f $(CERT) ] || [ ! -f $(KEY) ]; then \
		openssl req -x509 -newkey rsa:2048 \
			-keyout $(KEY) -out $(CERT) \
			-days 365 -nodes \
			-subj "/CN=localhost" 2>/dev/null; \
		echo "generated $(CERT) / $(KEY)"; \
	else \
		echo "certificate already exists"; \
	fi

test:
	go test ./... -count=1

vet:
	go vet ./...

bench:
	go test ./... -bench=. -benchtime=1s

clean:
	rm -f $(BIN) $(BIN_SRV) $(CERT) $(KEY)
	go clean ./...

help:
	@echo "Usage:"
	@echo "  make all        build binaries + generate dev certificate"
	@echo "  make build      build client and server binaries"
	@echo "  make client     build client binary only"
	@echo "  make server     build server binary only"
	@echo "  make gen-cert   generate self-signed TLS certificate"
	@echo "  make test       run all tests"
	@echo "  make bench      run all benchmarks"
	@echo "  make vet        run go vet"
	@echo "  make clean      remove binaries and certificates"
