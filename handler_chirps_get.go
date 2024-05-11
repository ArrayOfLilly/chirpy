package main

import (
	"net/http"
	"sort"
	"strconv"
)

// handlerChirpGet handles the retrieval of a specific chirp based on the provided chirp ID.
//
// It takes in an http.ResponseWriter and an http.Request as parameters.
// The function retrieves the chirp ID from the request path value, converts it to an integer, and checks for any conversion errors.
// If the chirp ID is invalid, it responds with a 400 Bad Request status code and an error message.
// If the chirp ID is valid, the function retrieves the corresponding chirp from the database using the provided DB instance.
// If the chirp cannot be found, it responds with a 404 Not Found status code and an error message.
// If the chirp is found, it responds with a 200 OK status code and a JSON representation of the chirp.
// The function does not return any values.
func (cfg *apiConfig) handlerChirpGet(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	chirpId, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	dbChirp, err := cfg.DB.GetChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}
	
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:   dbChirp.ID,
		Body: dbChirp.Body,
	})
}

// handlerChirpsRetrieve retrieves chirps from the database and responds with a JSON
// array of Chirp objects sorted by ID in ascending order. If there is an error
// retrieving the chirps, it responds with a 500 status code and an error message.
//
// Parameters:
// - w: http.ResponseWriter - the response writer used to write the response.
// - r: *http.Request - the HTTP request.
//
// Return type: None.
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}