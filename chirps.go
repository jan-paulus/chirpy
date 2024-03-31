package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	const maxChripLength = 140
	if len(params.Body) > maxChripLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleanedText := getCleanedText(params.Body, badWords)

	chirp, err := cfg.DB.CreateChirp(cleanedText)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

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

  chirpId, err := strconv.Atoi(idParam);
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

func getCleanedText(text string, badWords map[string]struct{}) string {
	words := strings.Split(text, " ")

	for i, word := range words {
		lword := strings.ToLower(word)
		if _, ok := badWords[lword]; ok {
			words[i] = "****"
		}
	}

	return strings.TrimSpace(strings.Join(words, " "))
}
