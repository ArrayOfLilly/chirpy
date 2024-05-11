package chirpy

import (
	"log"
	"net/http"

	"github.com/ArrayOfLilly/chirpy/internal/database"
)

// a struct that will hold any stateful, in-memory data we'll need to keep track of
	type apiConfig struct {
		DB *database.DB
		fileserverHits int
	}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// initializing apiConfig
	cfg := &apiConfig{

		fileserverHits: 0,
	}

	// ServeMux is an HTTP request multiplexer. 
	// It matches the URL of each incoming request against a list of registered patterns and 
	// calls the handler for the pattern that most closely matches the URL.
	mux := http.NewServeMux()

	// Handle add the handler (http.FileServer([path])) to the specified request ("/")
	// StripPrefix for set alternate route to the request
	// middlewareMetricsInc add the fileserverHit coumt functionality
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	// mux.HandleFunc("POST /api/chirps", handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", cfg.handlerChirpsRead)


	// readiness endpoint to shpow server status for external 
	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	// write and send the metrics
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)

	// reset the counter
	mux.HandleFunc("/api/reset", cfg.handlerReset)

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
