package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithError handles responding with an error message.
//
// Parameters:
// - w: http.ResponseWriter - the response writer used to write the response.
// - code: int - the HTTP status code for the error response.
// - msg: string - the error message to be included in the response.
// Return type: None.
func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

// respondWithJSON generates a JSON response with the provided status code and payload.
//
// Parameters:
// w http.ResponseWriter - the response writer to write the JSON response to.
// code int - the status code to return.
// payload interface{} - the data to be converted to JSON and written to the response.
// Return type(s): void
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}