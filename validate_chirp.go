package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	w.Header().Set("Content-Type", "application/json")

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

	respondWithJSON(w, http.StatusOK, successResponse{
		CleanedBody: getCleanedText(params.Body, badWords),
	})
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
