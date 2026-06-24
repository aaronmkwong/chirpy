package main

import (
	"encoding/json"
	"net/http"

	"github.com/aaronmkwong/chirpy/internal/auth"
)

// handlerLogin authenticates a user using email and password.
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// request body structure
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// decode request body
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request")
		return
	}

	// retrieve user by email
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// compare supplied password to stored password hash
	passwordMatches, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// reject invalid password
	if !passwordMatches {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// return authenticated user without password hash
	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}