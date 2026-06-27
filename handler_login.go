package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaronmkwong/chirpy/internal/auth"
)

// handlerLogin authenticates a user using email and password.
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	// Request body structure
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	// Response structure
	type loginResponse struct {
		User
		Token string `json:"token"`
	}

	// Decode request body
	params := parameters{}
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't decode request",
		)
		return
	}

	// Retrieve user by email
	user, err := cfg.db.GetUserByEmail(
		r.Context(),
		params.Email,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
		)
		return
	}

	// Compare supplied password against stored password hash
	passwordMatches, err := auth.CheckPasswordHash(
		params.Password,
		user.HashedPassword,
	)
	if err != nil || !passwordMatches {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
		)
		return
	}

	// Default JWT lifetime to one hour
	expiresIn := time.Hour

	// Use the requested expiration if it is positive and does not exceed one hour
	if params.ExpiresInSeconds > 0 &&
		params.ExpiresInSeconds <= int(time.Hour.Seconds()) {
		expiresIn = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	// Create JWT
	token, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expiresIn,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create token",
		)
		return
	}

	// Return authenticated user with JWT
	respondWithJSON(
		w,
		http.StatusOK,
		loginResponse{
			User: User{
				ID:        user.ID,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				Email:     user.Email,
			},
			Token: token,
		},
	)
}