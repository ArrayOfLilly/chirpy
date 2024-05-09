package main

import (
	"fmt"
	"log"
	"net/http"
)

// a struct that will hold any stateful, in-memory data we'll need to keep track of
	type apiConfig struct {
		fileserverHits int
	}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// initializing apiConfig
	apicfg := &apiConfig{
		fileserverHits: 0,
	}

	// ServeMux is an HTTP request multiplexer. 
	// It matches the URL of each incoming request against a list of registered patterns and 
	// calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()

	// Handle add the handler (http.FileServer([path])) to the specified request ("/")
	// StripPrefix for set alternate route to the request
	// middlewareMetricsInc add the fileserverHit coumt functionality
	mux.Handle("/app/*", apicfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// readiness endpoint to shpow server status for external 
	mux.HandleFunc("GET /healthz", handlerReadiness)

	// write and send the metrics
	mux.HandleFunc("GET /metrics", apicfg.handlerMetrics)

	// reset the counter
	mux.HandleFunc("/reset", apicfg.handlerReset)

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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(200)))
	}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
	}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
	}


