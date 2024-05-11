package chirpy

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	// request body json maps to this struct
	type parameters struct {
		// these tags indicate how the keys in the JSON should be mapped to the struct fields
        // the struct fields must be exported (start with a capital letter) if you want them parsed
		Body string `json:"body"`
	}

	// this struct maps to response body json
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// error handling while decoding request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// validate chirp, invalid case
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// valid case
	filter := []string{"kerfuffle", "sharbert", "fornax"}
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: filterChirp(filter, params.Body),
	})
}

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

func filterChirp(forbiddenWords []string, chirp string) string {
	forbiddenWordsDict := map[string]bool{}
	for _, w := range forbiddenWords {
		forbiddenWordsDict[w] = true
	}

	permittedWords := []string{}

	words := strings.Split(chirp, " ")
	for _, w := range words {
		_, ok := forbiddenWordsDict[strings.ToLower(w)]
		if ok {
			permittedWords = append(permittedWords, "****")
		} else {
			permittedWords = append(permittedWords, w)
		}
	}

	cleanedChirp := strings.Join(permittedWords, " ")
	return cleanedChirp
}
