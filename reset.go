package main

import "net/http"

// handlerReset handles the reset functionality of the API hit counter.
//
// It resets the value of `fileserverHits` in the `apiConfig` struct to 0.
// It writes the HTTP status code 200 (OK) and the response body "Hits reset to 0" to the `http.ResponseWriter`.
//
// Parameters:
// - w: http.ResponseWriter - the response writer used to write the response.
// - r: *http.Request - the HTTP request.
//
// Return type: None.
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}