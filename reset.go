package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) resetHandler(rw http.ResponseWriter, rq *http.Request) {
	cfg.fileserverHits.Store(0)
	log.Printf("Hits reset to 0")
	rw.WriteHeader(http.StatusOK)
}
