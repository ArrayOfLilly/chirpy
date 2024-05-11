package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ArrayOfLilly/chirpy/internal/auth"
	"github.com/ArrayOfLilly/chirpy/internal/database"
)

func (cfg *apiConfig) handlerAutentication(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if  err != nil {
		if !errors.Is(err, database.ErrNotExist){
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithError(w, http.StatusBadRequest, "User does not exists")
			return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		Email: user.Email,
	})
}
