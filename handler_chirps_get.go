package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, _ *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		log.Printf("Error fetching chirps from db: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
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
