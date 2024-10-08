package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port = "8080"

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new apiConfig instance
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Register handler functions
	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthCheckHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)

	// Create new server instance
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Run the server
	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
