package main

import (
	"net/http"
	"github.com/google/uuid"
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

// handlerChirpGet retrieves a single chirp by ID and returns it as JSON.
func (cfg *apiConfig) handlerChirpGet(
w http.ResponseWriter,
r *http.Request,
) {
// Get the chirp ID from the URL path parameter
chirpIDString := r.PathValue("chirpID")

// Convert the string ID to a UUID
chirpID, err := uuid.Parse(chirpIDString)
if err != nil {
	respondWithError(
		w,
		http.StatusNotFound,
		"Chirp not found",
	)
	return
}

// Query the chirp from the database
dbChirp, err := cfg.db.GetChirp(
	r.Context(),
	chirpID,
)
if err != nil {
	respondWithError(
		w,
		http.StatusNotFound,
		"Chirp not found",
	)
	return
}

// Map the database model to our API-facing structure
chirp := Chirp{
	ID:        dbChirp.ID,
	CreatedAt: dbChirp.CreatedAt,
	UpdatedAt: dbChirp.UpdatedAt,
	Body:      dbChirp.Body,
	UserID:    dbChirp.UserID,
}

// Return the chirp as JSON
respondWithJSON(
	w,
	http.StatusOK,
	chirp,
)

}
