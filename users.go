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

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type CreateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) createUserHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := CreateUserParams{}
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

	// Return error if password property is missing or empty
	if len(params.Password) == 0 {
		err = respondWithError(rw, http.StatusBadRequest, "Password is required")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Hash user password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error hashing password")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Insert user into database
	dbParams := database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email:          params.Email,
	}

	dbUser, err := cfg.db.CreateUser(rq.Context(), dbParams)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error creating user")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Map database user to User struct
	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	// Return user
	err = respondWithJSON(rw, http.StatusCreated, user)
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}

func (cfg *apiConfig) loginHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := LoginParams{}
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

	// Return error if password property is missing or empty
	if len(params.Password) == 0 {
		err = respondWithError(rw, http.StatusBadRequest, "Password is required")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Get user from database
	dbUser, err := cfg.db.GetUserByEmail(rq.Context(), params.Email)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Not authorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Compare user password with hashed password
	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Not authorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Create access token
	token, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, time.Second*time.Duration(3600))
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error creating token")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Create refresh token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error creating refresh token")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Insert refresh token into database
	dbParams := database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	}

	dbRefreshToken, err := cfg.db.CreateRefreshToken(rq.Context(), dbParams)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error creating refresh token")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Map database user to User struct
	user := User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: dbRefreshToken.Token,
		IsChirpyRed:  dbUser.IsChirpyRed,
	}

	// Return user
	err = respondWithJSON(rw, http.StatusOK, user)
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}

func (cfg *apiConfig) updateUserHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := UpdateUserParams{}
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

	// Return error if email property is missing or empty
	if len(params.Email) == 0 {
		err = respondWithError(rw, http.StatusBadRequest, "Email is required")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Return error if password property is missing or empty
	if len(params.Password) == 0 {
		err = respondWithError(rw, http.StatusBadRequest, "Password is required")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Hash user password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error hashing password")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Update user in database
	dbParams := database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	dbUser, err := cfg.db.UpdateUser(rq.Context(), dbParams)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error updating user")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Map database user to User struct
	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	// Return user
	err = respondWithJSON(rw, http.StatusOK, user)
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}
