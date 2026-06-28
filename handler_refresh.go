package main

import (
	"net/http"
	"time"

	"github.com/aaronmkwong/chirpy/internal/auth"
)

// handlerRefresh exchanges a valid refresh token for a new access token.
func (cfg *apiConfig) handlerRefresh(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Response structure
	type response struct {
		Token string `json:"token"`
	}

	// Extract the refresh token from the Authorization header.
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Missing or invalid refresh token",
		)
		return
	}

	// Retrieve the associated user.
	user, err := cfg.db.GetUserFromRefreshToken(
		r.Context(),
		refreshToken,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Invalid refresh token",
		)
		return
	}

	// Create a new one-hour access token.
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

	// Return the new access token.
	respondWithJSON(
		w,
		http.StatusOK,
		response{
			Token: token,
		},
	)
}