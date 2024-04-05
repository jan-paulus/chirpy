package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jan-paulus/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s", err)
	}

	cfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}

	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/*", cfg.middlewareMetricsInc(fsHandler))

	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
  mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
  
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/login", cfg.handlerLoginUser)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("/api/reset", cfg.handlerMetricsReset)

	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
