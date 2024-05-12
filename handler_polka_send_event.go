package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ArrayOfLilly/chirpy/internal/auth"
	"github.com/ArrayOfLilly/chirpy/internal/database"
)

func (cfg *apiConfig) handlerPolkaSendEvent(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	type response struct {}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find API key")
		return
	}
	
	if !auth.CheckApiKey(cfg.polkaApiKey, apiKey) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized user")
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "")
		return
	}

	user, err := cfg.DB.GetUser(params.Data.UserId)
		if err != nil {
			if errors.Is(database.ErrNotExist, err) {
				respondWithError(w, http.StatusNotFound, "Couldn't find user ID")
				return
			}
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

	err = cfg.DB.UpdateUserIsChirpyRed(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, response{})
}