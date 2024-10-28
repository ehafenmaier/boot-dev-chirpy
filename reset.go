package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(rw http.ResponseWriter, rq *http.Request) {
	// Return forbidden if platform is not "dev"
	if cfg.platform != "dev" {
		err := respondWithError(rw, http.StatusForbidden, "Forbidden")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	// Reset users table
	err := cfg.db.ResetUsers(rq.Context())
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error resetting users")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	cfg.fileserverHits.Store(0)
	log.Printf("Hits reset to 0")
	rw.WriteHeader(http.StatusOK)
}
