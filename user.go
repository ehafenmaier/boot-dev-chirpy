package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type createUserParams struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) createUserHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := createUserParams{}
	err := decoder.Decode(&params)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Invalid request payload")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Return error if email property is missing or empty
	if len(params.Email) == 0 {
		err = respondWithError(rw, http.StatusBadRequest, "Email is required")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Insert user into database
	dbUser, err := cfg.db.CreateUser(rq.Context(), params.Email)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error creating user")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Map database user to User struct
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	// Return user
	data, err := json.Marshal(user)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error marshalling response")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	rw.WriteHeader(http.StatusCreated)
	_, err = rw.Write(data)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error writing response")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}
}
