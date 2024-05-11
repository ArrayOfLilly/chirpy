package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"

	"github.com/ArrayOfLilly/chirpy/internal/auth"
	"github.com/ArrayOfLilly/chirpy/internal/database"
)

type User struct {
	ID   int    `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}


func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	// represents the parameters expected in the request body.
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

	validEmail, err := mail.ParseAddress(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := cfg.DB.CreateUser(string(validEmail.Address), hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:   user.ID,
		Email: user.Email,
	})
}