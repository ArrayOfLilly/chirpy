package main

import (
	"fmt"
	"net/http"
)

// handlerMetrics generates an HTML response containing metrics about the Chirpy /app/* usage.
//
// Parameters:
// - w http.ResponseWriter: the response writer to write the HTML response to.
// - r *http.Request: the HTTP request.
//
// Return type: None.
func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %dtimes!</p>
</body>

</html>
	`, cfg.fileserverHits)))
}

// middlewareMetricsInc extends the next http.Handler functionality, 
// by incrementing the number of fileserver hits (on /app/*) in the apiConfig struct,
// (injecting the corresponding code into it).
//
// Parameters:
// - next http.Handler: the next http.Handler in the middleware chain.
//
// Return type: http.Handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}