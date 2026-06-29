package main

import (
	"encoding/json"
	"net/http"

	"github.com/aaronmkwong/chirpy/internal/auth"
	"github.com/aaronmkwong/chirpy/internal/database"
)

// handlerUsersUpdate updates the authenticated user's email and password.
func (cfg *apiConfig) handlerUserUpdate(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Request body structure.
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the request body.
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

	// Extract the bearer token.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Missing or invalid token",
		)
		return
	}

	// Validate the JWT and obtain the authenticated user's ID.
	userID, err := auth.ValidateJWT(
		token,
		cfg.jwtSecret,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Missing or invalid token",
		)
		return
	}

	// Hash the new password.
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't hash password",
		)
		return
	}

	// Update the user.
	user, err := cfg.db.UpdateUser(
		r.Context(),
		database.UpdateUserParams{
			ID:             userID,
			Email:          params.Email,
			HashedPassword: hashedPassword,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't update user",
		)
		return
	}

	// Return the updated user (without the password hash).
	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	)
}