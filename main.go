package main

import (
	"net/http"
	"log"
	)

func main() {
	const port = "8080"

	// ServeMux is an HTTP request multiplexer. 
	// It matches the URL of each incoming request against a list of registered patterns and 
	// calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()

	// A Server defines parameters for running an HTTP server. 
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	// logs: 2024/05/09 15:00:03 Serving on port: 8080

	// ListenAndServe listens on the TCP network address srv.Addr and 
	// then calls Serve to handle requests on incoming connections.
	// Serve accepts incoming connections on the Listener l, 
	// creating a new service goroutine for each. 
	// The service goroutines read requests and then call srv.Handler to reply to them.
	log.Fatal(srv.ListenAndServe())
}

