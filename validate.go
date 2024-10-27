package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type validateChirpParams struct {
	Body string `json:"body"`
}

type returnCleaned struct {
	CleanedBody string `json:"cleaned_body"`
}

func validateChirpHandler(rw http.ResponseWriter, rq *http.Request) {
	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := validateChirpParams{}
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

	// Clean the chirp
	respBody := returnCleaned{
		CleanedBody: replaceBadWords(params.Body),
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error marshalling response")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)
	if err != nil {
		err = respondWithError(rw, http.StatusInternalServerError, "Error writing response")
		if err != nil {
			log.Printf("Error responding: %v", err)
		}
		return
	}
}
