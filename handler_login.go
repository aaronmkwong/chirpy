package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aaronmkwong/chirpy/internal/auth"
	"github.com/aaronmkwong/chirpy/internal/database"
)

// handlerLogin authenticates a user and returns both an access token and a refresh token.
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	// Request body structure
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	// Response structure
	type loginResponse struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	// Create JWT access token (always expires after one hour)
	token, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create access token",
		)
		return
	}

	// Generate a refresh token
	refreshToken := auth.MakeRefreshToken()

	// Store the refresh token in the database
	_, err = cfg.db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:      refreshToken,
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
			UserID:     user.ID,
			ExpiresAt:  time.Now().UTC().Add(60 * 24 * time.Hour),
			RevokedAt:  sql.NullTime{},
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create refresh token",
		)
		return
	}

	// Return authenticated user with both tokens
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
			Token:        token,
			RefreshToken: refreshToken,
		},
	)
}