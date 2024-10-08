package main

import (
	"fmt"
	"log"
	"net/http"
)

func (cfg *apiConfig) hitsHandler(rw http.ResponseWriter, rq *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.WriteHeader(http.StatusOK)

	body := "<html><body><h1>Welcome, Chirpy Admin</h1><p>" +
		"Chirpy has been visited %d times!" +
		"</p></body></html>"

	_, err := rw.Write([]byte(fmt.Sprintf(body, cfg.fileserverHits.Load())))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// Middleware function to increment the fileserverHits counter
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf("Hits: %d", cfg.fileserverHits.Load())
		next.ServeHTTP(rw, rq)
	})
}
