package main

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

// Structure of a "Chirp" (limited length message)
type User struct {
	ID   int    `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

// handlerUsersCreate handles the creation of a new user based on the request body.
//
// It decodes the request body into parameters, validates the email address, and creates a user in the database.
// It returns the created user in the response.
//
// Parameters:
// - w http.ResponseWriter: the response writer to write the HTTP response to.
// - r *http.Request: the HTTP request.
//
// Return type: None.
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 4)
	if err != nil {
		respondWithError(w, http.StatusBadRequest,  err.Error())
		return
	}
	
	user, err := cfg.DB.CreateUser(string(validEmail.Address), string(passwordHash))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:   user.ID,
		Email: user.Email,
	})
}