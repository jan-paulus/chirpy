package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jan-paulus/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	responseBody := response{
		Email: user.Email,
		Id:    user.Id,
	}
	respondWithJSON(w, http.StatusCreated, responseBody)
}
