package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register handler functions
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthCheckHandler)

	// Create new server instance
	srv := &http.Server {
		Addr:    ":" + port,
		Handler: mux,
	}

	// Run the server
	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}

// Handler function for health check endpoint
func healthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}