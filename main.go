package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the handler function for the root URL pattern
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Create new server instance
	srv := &http.Server {
		Addr:    ":" + port,
		Handler: mux,
	}

	// Run the server
	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}