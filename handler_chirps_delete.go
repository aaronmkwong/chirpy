package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/aaronmkwong/chirpy/internal/auth"
	"github.com/google/uuid"
)

// handlerChirpsDelete deletes a chirp if it belongs to the authenticated user.
func (cfg *apiConfig) handlerChirpDelete(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Extract the bearer token from the Authorization header.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Missing or invalid token",
		)
		return
	}

	// Validate the JWT and retrieve the authenticated user's ID.
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

	// Parse the chirp ID from the URL path.
	chirpID, err := uuid.Parse(
		r.PathValue("chirpID"),
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Invalid chirp ID",
		)
		return
	}

	// Retrieve the chirp from the database.
	chirp, err := cfg.db.GetChirp(
		r.Context(),
		chirpID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(
				w,
				http.StatusNotFound,
				"Chirp not found",
			)
			return
		}

		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't retrieve chirp",
		)
		return
	}

	// Ensure the authenticated user owns the chirp.
	if chirp.UserID != userID {
		respondWithError(
			w,
			http.StatusForbidden,
			"Forbidden",
		)
		return
	}

	// Delete the chirp.
	err = cfg.db.DeleteChirp(
		r.Context(),
		chirpID,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't delete chirp",
		)
		return
	}

	// Successfully deleted.
	w.WriteHeader(http.StatusNoContent)
}