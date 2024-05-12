package main

import (
	"net/http"
	"sort"
	"strconv"
)

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
		AuthorID: dbChirp.AuthorID,
	})
}

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	authorIdStr := r.URL.Query().Get("author_id")
	if authorIdStr != ""{
		authorId, err := strconv.Atoi(authorIdStr)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse authorID ID")
			return
		}
		chirps := []Chirp{}
		for _, dbChirp := range dbChirps {
			if dbChirp.AuthorID == authorId {
				chirps = append(chirps, Chirp{
				ID:   dbChirp.ID,
				Body: dbChirp.Body,
				AuthorID: dbChirp.AuthorID,
				})
			}
		}
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}