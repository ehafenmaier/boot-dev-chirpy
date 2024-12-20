package main

import (
	"encoding/json"
	"github.com/ehafenmaier/boot-dev-chirpy/internal/auth"
	"github.com/ehafenmaier/boot-dev-chirpy/internal/database"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

type createChirpParams struct {
	Body string `json:"body"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(rw http.ResponseWriter, rq *http.Request) {
	// Decode request body
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

	// Validate bearer token
	token, err := auth.GetBearerToken(rq.Header)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Validate JWT
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
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
		UserID: userID,
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

func (cfg *apiConfig) getChirpsHandler(rw http.ResponseWriter, rq *http.Request) {
	// Check for sort parameter
	sortParam := rq.URL.Query().Get("sort")
	if len(sortParam) == 0 {
		sortParam = "asc"
	}

	if sortParam != "asc" && sortParam != "desc" {
		err := respondWithError(rw, http.StatusBadRequest, "Invalid sort parameter")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Check for author id query parameter
	authorParam := rq.URL.Query().Get("author_id")
	if len(authorParam) > 0 {
		authorID, err := uuid.Parse(authorParam)
		if err != nil {
			err = respondWithError(rw, http.StatusBadRequest, "Invalid author ID")
			if err != nil {
				log.Printf("Error responding: %v", err)
			}
			return
		}

		// Get chirps by author from database
		dbParams := database.GetChirpsByUserIdParams{
			UserID:  authorID,
			Column2: sortParam,
		}
		dbChirps, err := cfg.db.GetChirpsByUserId(rq.Context(), dbParams)
		if err != nil {
			err = respondWithError(rw, http.StatusInternalServerError, "Error getting chirps")
			if err != nil {
				log.Printf("Error responding: %v", err)
			}
			return
		}

		// Map database chirps to Chirp struct
		chirps := make([]Chirp, len(dbChirps))
		for i, dbChirp := range dbChirps {
			chirps[i] = Chirp{
				ID:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				Body:      dbChirp.Body,
				UserID:    dbChirp.UserID,
			}
		}

		// Return chirps
		err = respondWithJSON(rw, http.StatusOK, chirps)
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Get chirps from database
	dbChirps, err := cfg.db.GetChirps(rq.Context(), sortParam)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error getting chirps")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Map database chirps to Chirp struct
	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	// Return chirps
	err = respondWithJSON(rw, http.StatusOK, chirps)
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}

func (cfg *apiConfig) getChirpHandler(rw http.ResponseWriter, rq *http.Request) {
	// Get chirp ID from URL
	id, err := uuid.Parse(rq.PathValue("id"))
	if err != nil {
		err = respondWithError(rw, http.StatusBadRequest, "Invalid chirp ID")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Get chirp from database
	dbChirp, err := cfg.db.GetChirpById(rq.Context(), id)
	if err != nil {
		err = respondWithError(rw, http.StatusNotFound, "Chirp not found")
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
	err = respondWithJSON(rw, http.StatusOK, chirp)
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}

func (cfg *apiConfig) deleteChirpHandler(rw http.ResponseWriter, rq *http.Request) {
	// Get chirp ID from URL
	id, err := uuid.Parse(rq.PathValue("id"))
	if err != nil {
		err = respondWithError(rw, http.StatusBadRequest, "Invalid chirp ID")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Validate bearer token
	token, err := auth.GetBearerToken(rq.Header)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Validate JWT
	userID, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Get chirp from database
	dbChirp, err := cfg.db.GetChirpById(rq.Context(), id)
	if err != nil {
		err = respondWithError(rw, http.StatusNotFound, "Chirp not found")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Return error if user is not chirp owner
	if dbChirp.UserID != userID {
		err = respondWithError(rw, http.StatusForbidden, "Forbidden")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Delete chirp from database
	err = cfg.db.DeleteChirp(rq.Context(), id)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error deleting chirp")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Return success (no content)
	respondWithNoContent(rw)
}
