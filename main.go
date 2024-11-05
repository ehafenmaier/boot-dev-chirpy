package main

import (
	"database/sql"
	"github.com/ehafenmaier/boot-dev-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
}

func main() {
	const port = "8080"

	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection string
	dbURL := os.Getenv("DB_URL")

	// Open connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a new apiConfig instance
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       os.Getenv("PLATFORM"),
		tokenSecret:    os.Getenv("TOKEN_SECRET"),
	}

	// Register handler functions
	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthCheckHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.hitsHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.getChirpHandler)
	mux.HandleFunc("DELETE /api/chirps/{id}", cfg.deleteChirpHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.tokenRefreshHandler)
	mux.HandleFunc("POST /api/revoke", cfg.tokenRevokeHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateUserHandler)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.polkaWebhookHandler)

	// Create new server instance
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Run the server
	log.Printf("Server started on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
