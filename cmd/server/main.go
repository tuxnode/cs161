package main

import (
	"flag"
	"log"

	"github.com/cs161-staff/project2-starter-code/internal/server"
)

func main() {
	addr := flag.String("address", "localhost:8080", "listen address")
	dir := flag.String("dir", "./data", "data directory")
	cert := flag.String("cert", "", "TLS cert file (required if -tls=true)")
	key := flag.String("key", "", "TLS key file (required if -tls=true)")
	tlsEnabled := flag.Bool("tls", true, "enable TLS encryption")
	flag.Parse()

	if *tlsEnabled && (*cert == "" || *key == "") {
		log.Fatal("-cert and -key are required when -tls=true")
	}

	srv, err := server.New(server.Config{
		Addr:       *addr,
		DataDir:    *dir,
		Cert:       *cert,
		Key:        *key,
		TLSEnabled: *tlsEnabled,
	})
	if err != nil {
		log.Fatalf("server init: %v", err)
	}

	log.Fatal(srv.Run())
}
