package main

import (
	"encoding/json"
	"log"
	"net/http"
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
