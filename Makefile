BIN     := sharelock
BIN_SRV := sharelock-server
CERT    := cert.pem
KEY     := key.pem

.PHONY: all build client server clean test vet bench gen-cert

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
