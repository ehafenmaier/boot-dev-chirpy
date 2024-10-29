package main

import (
	"encoding/json"
	"github.com/ehafenmaier/boot-dev-chirpy/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type createChirpParams struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := createChirpParams{}
	err := decoder.Decode(&params)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Invalid request payload")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Return error if chirp longer than 140 characters
	if len(params.Body) > 140 {
		err = respondWithError(rw, http.StatusBadRequest, "Chirp is too long")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Insert chirp into database
	dbParams := database.CreateChirpParams{
		Body:   replaceBadWords(params.Body),
		UserID: params.UserID,
	}

	dbChirp, err := cfg.db.CreateChirp(rq.Context(), dbParams)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error creating chirp")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Map database chirp to Chirp struct
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	// Return chirp
	err = respondWithJSON(rw, http.StatusCreated, chirp)
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}
