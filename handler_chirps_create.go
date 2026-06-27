package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aaronmkwong/chirpy/internal/auth"
	"github.com/aaronmkwong/chirpy/internal/database"
	"github.com/google/uuid"
)

// Chirp defines the API response structure for a chirp.
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// handlerChirpsCreate handles POST /api/chirps. It authenticates the user,
// validates the chirp, stores it in the database, and returns the created resource.
func (cfg *apiConfig) handlerChirpsCreate(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Request body structure.
	type parameters struct {
		Body string `json:"body"`
	}

	params := parameters{}

	// Decode the request body.
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

	// Extract the bearer token from the Authorization header.
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Unauthorized",
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
			"Unauthorized",
		)
		return
	}

	// Validate and clean the chirp.
	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	// Persist the chirp.
	dbChirp, err := cfg.db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: userID,
		},
	)
	if err != nil {
		log.Printf("Error creating chirp: %v", err)
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create chirp",
		)
		return
	}

	// Convert the database model to the API response model.
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	// Return the created chirp.
	respondWithJSON(
		w,
		http.StatusCreated,
		chirp,
	)
}

// validateChirp enforces the 140-character limit and censors forbidden words.
func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", errors.New("Chirp is too long")
	}

	return getCleanedChirp(chirp), nil
}

// getCleanedChirp replaces forbidden words with asterisks.
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
