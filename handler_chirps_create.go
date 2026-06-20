package main

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"net/http"
	"strings"

	"github.com/aaronmkwong/chirpy/internal/database"
	"github.com/google/uuid"
)

// define chirp response 
type Chirp struct {
    ID        uuid.UUID `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Body      string    `json:"body"`
    UserID    uuid.UUID `json:"user_id"`
}

// handlerChirpsCreate handles POST /api/chirps: it validates the incoming
// chirp, saves it to the database, and returns the created resource.
func (cfg *apiConfig) handlerChirpsCreate(
	w http.ResponseWriter,
	r *http.Request,
) {
	// parameters describes the shape of the expected JSON request body.
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{}

	// Decode the request body into params.
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		// Malformed JSON: nothing more we can do, so report a server error.
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Validate length and scrub profanity, returning the cleaned body.
	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		// errorResponse is the JSON shape returned when validation fails.
		type errorResponse struct {
			Error string `json:"error"`
		}

		respBody := errorResponse{
			Error: err.Error(),
		}

		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Validation failure is the client's fault: respond with 400.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	// Persist the chirp. The database generates id, created_at, and updated_at.
	chirp, err := cfg.db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: params.UserID,
		},
	)
	if err != nil {
		// Insert failed (e.g. unknown user_id or DB issue).
		log.Printf("Error creating chirp: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Map the database chirp into our tagged response type, then serialize.
	dat, err := json.Marshal(Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success: respond with 201 Created and the new chirp.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(dat)
}

// validateChirp enforces the 140-character limit. If the chirp is valid,
// it returns the profanity-scrubbed body; otherwise it returns an error.
func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("Chirp is too long")
	}
	return getCleanedChirp(chirp), nil
}

// getCleanedChirp replaces forbidden words with asterisks and returns
// the cleaned string. It never fails.
func getCleanedChirp(chirp string) string {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(chirp, " ")

	for i, word := range words {
		if badWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}