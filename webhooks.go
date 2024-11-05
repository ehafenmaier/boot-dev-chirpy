package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type PolkaData struct {
	UserID uuid.UUID `json:"user_id"`
}

type PolkaWebhookParams struct {
	Event string    `json:"event"`
	Data  PolkaData `json:"data"`
}

func (cfg *apiConfig) polkaWebhookHandler(rw http.ResponseWriter, rq *http.Request) {
	// Decode request body
	decoder := json.NewDecoder(rq.Body)
	params := PolkaWebhookParams{}
	err := decoder.Decode(&params)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Invalid request payload")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Check for user upgraded event
	if params.Event != "user.upgraded" {
		respondWithNoContent(rw)
		return
	}

	// Get user from database
	dbUser, err := cfg.db.GetUserByID(rq.Context(), params.Data.UserID)
	if err != nil {
		err = respondWithError(rw, http.StatusNotFound, "User not found")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Upgrade user to chirpy red
	dbUser, err = cfg.db.UpgradeUserToChirpyRed(rq.Context(), dbUser.ID)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error upgrading user to chirpy red")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Return success (no content)
	respondWithNoContent(rw)
}
