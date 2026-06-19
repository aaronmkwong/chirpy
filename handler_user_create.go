package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log"

	"github.com/google/uuid"
)

// Request body structure
type requestParameters struct {
	Email string `json:"email"`
}

// Internal User structure to control JSON output format
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	// decode the request body JSON
	params := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// call the SQLC-generated CreateUser database method
	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
		return
	}

	// map the returned database.User to the local User struct
	responseUser := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	// respond with 201 Created and the user as JSON
	respondWithJSON(w, http.StatusCreated, responseUser)
}
