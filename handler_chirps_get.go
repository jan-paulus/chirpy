package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/jan-paulus/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Printf("Error fetching chirps from db: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if s == "" {
		respondWithJSON(w, http.StatusOK, chirps)
	}

	authorId, err := strconv.Atoi(s)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	filteredChirps := []database.Chirp{}

	for _, chrip := range chirps {
		if chrip.AuthorId == authorId {
			filteredChirps = append(filteredChirps, chrip)
		}
	}

	respondWithJSON(w, http.StatusOK, filteredChirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("chirpID")
	if idParam == "" {
		respondWithError(w, http.StatusBadRequest, "Please provide an id")
	}

	chirpId, err := strconv.Atoi(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to convert ID to int")
		return
	}

	chirp, err := cfg.DB.GetChirp(chirpId)
	if err != nil {

		if err.Error() == "Chirp does not exist" {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		log.Printf("Error fetching chirp from db: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
