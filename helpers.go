package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type returnError struct {
	Error string `json:"error"`
}

func respondWithJSON(rw http.ResponseWriter, code int, payload interface{}) error {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	rw.WriteHeader(code)
	_, err = rw.Write(response)
	if err != nil {
		return err
	}

	return nil
}

func respondWithError(rw http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(rw, code, returnError{Error: msg})
}

func respondWithNoContent(rw http.ResponseWriter) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusNoContent)
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
