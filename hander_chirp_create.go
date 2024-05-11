package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Structure of a "Chirp" (limited length message)
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}


// handlerChirpsCreate handles the creation of a new chirp based on the request body.
//
// It decodes the request body into parameters, validates the chirp content, and creates a chirp in the database.
// It returns the created chirp in the response.
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	// represents the parameters expected in the request body.
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(cleaned)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:   chirp.ID,
		Body: chirp.Body,
	})
}

// validateChirp validates the given chirp body.
//
// It checks if the length of the body is greater than the maximum allowed length (140 characters).
// If the length is greater, it returns an error indicating that the chirp is too long.
// It also checks if the body contains any bad words defined in the `badWords` map.
// If a bad word is found, it replaces it with asterisks.
// Finally, it returns the cleaned body and a nil error.
//
// Parameters:
// - body: the chirp body to be validated (string)
//
// Returns:
// - cleaned: the cleaned chirp body (string)
// - error: an error if the chirp is too long or contains bad words (error)
func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(body, badWords)
	return cleaned, nil
}

// getCleanedBody cleans the given body string by replacing any bad words with asterisks.
//
// Parameters:
// - body: the string containing the body to be cleaned (string)
// - badWords: a map of bad words to be replaced with asterisks (map[string]struct{})
//
// Returns:
// - cleaned: the cleaned body string with bad words replaced (string)
func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}