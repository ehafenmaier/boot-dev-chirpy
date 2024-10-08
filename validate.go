package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirpHandler(rw http.ResponseWriter, rq *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnError struct {
		Error string `json:"error"`
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}

	// Set response content type
	rw.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(rq.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding body: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return error if chirp longer than 140 characters
	if len(params.Body) > 140 {
		respBody := returnError{
			Error: "Chirp is too long",
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
		_, err = rw.Write(data)
		if err != nil {
			log.Printf("Error writing response: %s", err)
		}
		return
	}

	// Chirp is valid
	respBody := returnValid{
		Valid: true,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
	_, err = rw.Write(data)
	if err != nil {
		log.Printf("Error writing response: %s", err)
	}
}
