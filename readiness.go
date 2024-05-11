package main

import "net/http"

// handlerReadiness handles the readiness check for the application.
// Indicates the server state for external users
//
// It sets the Content-Type header to "text/plain; charset=utf-8" and writes the HTTP status code 200 (OK) along with the status text "OK" as the response body.
//
// Parameters:
// - w http.ResponseWriter: the response writer used to write the response.
// - r *http.Request: the HTTP request.
//
// Return type: None.
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}