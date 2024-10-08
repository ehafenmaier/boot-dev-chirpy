package main

import (
	"log"
	"net/http"
)

// Handler function for health check endpoint
func healthCheckHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	_, err := rw.Write([]byte(http.StatusText(http.StatusOK)))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
