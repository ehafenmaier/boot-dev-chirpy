package main

import (
	"github.com/ehafenmaier/boot-dev-chirpy/internal/auth"
	"log"
	"net/http"
	"time"
)

type RefreshedAccessToken struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) tokenRefreshHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	// Validate bearer token
	token, err := auth.GetBearerToken(rq.Header)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Get refresh token from database
	refreshToken, err := cfg.db.GetRefreshToken(rq.Context(), token)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Check if refresh token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Check if refresh token is revoked
	if refreshToken.RevokedAt.Valid {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Create refreshed access token
	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.tokenSecret, time.Second*time.Duration(3600))
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Respond with refreshed access token
	err = respondWithJSON(rw, http.StatusOK, RefreshedAccessToken{Token: accessToken})
	if err != nil {
		log.Printf("Error responding: %v", err)
	}
}

func (cfg *apiConfig) tokenRevokeHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	// Validate bearer token
	token, err := auth.GetBearerToken(rq.Header)
	if err != nil {
		err = respondWithError(rw, http.StatusUnauthorized, "Unauthorized")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Revoke refresh token
	err = cfg.db.RevokeRefreshToken(rq.Context(), token)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Internal Server Error")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Respond with success (no content)
	respondWithNoContent(rw)
}
