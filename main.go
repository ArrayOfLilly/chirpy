package main

import (
	"log"
	"net/http"

	"github.com/ArrayOfLilly/chirpy/internal/database"
)

// configuration object for the API
type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

// main is the entry point of the Go program.
//
// It initializes the necessary configurations and sets up the HTTP server.
// The function reads the database file "database.json" and creates a new database connection.
// It creates an instance of the apiConfig struct with the initial number of fileserver hits and the database connection.
// It creates a new HTTP request multiplexer and sets up the file server handler with the middleware for metrics increment.
// The file server handler is registered to handle requests with the "/app/*" pattern.
// The function registers various API handlers with the multiplexer.
// The handlerReadiness function is registered to handle "GET /api/healthz" requests.
// The handlerReset function of the apiConfig struct is registered to handle "GET /api/reset" requests.
// The handlerChirpsCreate function of the apiConfig struct is registered to handle "POST /api/chirps" requests.
// The handlerChirpsRetrieve function of the apiConfig struct is registered to handle "GET /api/chirps" requests.
// The handlerMetrics function of the apiConfig struct is registered to handle "GET /admin/metrics" requests.
// The function creates an HTTP server with the specified address and the multiplexer as the handler.
// It logs the file server root and the listening port.
// Finally, it starts the server and logs any errors that occur.
func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

