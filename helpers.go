package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type returnError struct {
	Error string `json:"error"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if payload == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		return err
	}
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, returnError{Error: msg})
}

func replaceBadWords(body string) string {
	// Bad words map
	badWords := map[string]string{
		"kerfuffle": "kerfuffle",
		"sharbert":  "sharbert",
		"fornax":    "fornax",
	}

	// Replace bad words
	bodySplit := strings.Split(body, " ")
	for i, word := range bodySplit {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			bodySplit[i] = "****"
		}
	}

	return strings.Join(bodySplit, " ")
}
