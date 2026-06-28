package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/aaronmkwong/chirpy/internal/auth"
	"github.com/aaronmkwong/chirpy/internal/database"
)

// handlerRevoke revokes a refresh token.
func (cfg *apiConfig) handlerRevoke(
	w http.ResponseWriter,
	r *http.Request,
) {

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

	// Capture the current UTC timestamp.
	now := time.Now().UTC()

	// Revoke the refresh token.
	err = cfg.db.RevokeRefreshToken(
		r.Context(),
		database.RevokeRefreshTokenParams{
			Token:     refreshToken,
			RevokedAt: sql.NullTime{Time: now, Valid: true},
			UpdatedAt: now,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't revoke refresh token",
		)
		return
	}

	// Return 204 No Content.
	w.WriteHeader(http.StatusNoContent)
}