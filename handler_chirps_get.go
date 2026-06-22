package main

import (
	"net/http"
)

// handlerChirpsGet retrieves all chirps from the database, sorts them,
// and returns them as a JSON array.
func (cfg *apiConfig) handlerChirpsGet(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Query all chirps from the database using our generated SQLC query
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't retrieve chirps",
		)
		return
	}

	// Slice to hold the API-formatted chirps, pre-allocated for efficiency
	chirps := make([]Chirp, 0, len(dbChirps))

	// Map the database-specific models to our public API-facing structures
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	// Send the mapped slice back to the client as a JSON response
	respondWithJSON(w, http.StatusOK, chirps)
}