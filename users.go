package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	user, err := cfg.DB.CreateUser(params.Email, params.Password)
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

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
    Token string `json:"token"`
	}
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "The email and password do not match with an existing user.")
		return
	}

	expiresAt := &jwt.NumericDate{}

	if params.ExpiresInSeconds != 0 {
		expiresAt = jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(params.ExpiresInSeconds)))
	} else {
		expiresAt = jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24))
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(user.Id),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: expiresAt,
  }
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  signedString, err := token.SignedString([]byte(cfg.jwtSecret))
  
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to authorize.")
		return
	}

	responseBody := response{
		Email: user.Email,
		Id:    user.Id,
    Token: signedString,
	}
	respondWithJSON(w, http.StatusOK, responseBody)
}
