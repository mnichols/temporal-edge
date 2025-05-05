package main

import (
	"fmt"
	//"github.com/mnichols/temporal-edge/generated/tmprl/v1"
	"github.com/mnichols/temporal-edge/tmprl/generated/tmprl/v1/tmprlv1connect"
	"github.com/mnichols/temporal-edge/tmprl/internal"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// The generated constructors return a path and a plain net/http
	// handler.
	mux.Handle(tmprlv1connect.NewTmprlServiceHandler(&internal.TmprlService{}))
	// Example Request
	/*
		curl http://localhost:8080/snailforce.v1.SnailforceService/Register \
			--header 'content-type: application/json' \
			--header 'accept: application/json \
			--data '{"id":"foo","value":"bar"}'
	*/
	fmt.Println("Starting Tmprl @", "localhost:8081")
	err := http.ListenAndServe(
		"localhost:8081",
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	log.Fatalf("listen failed: %v", err)
}
