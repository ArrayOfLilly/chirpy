package main

import (
	"net/http"
	"sort"
)

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