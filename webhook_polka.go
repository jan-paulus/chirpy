package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jan-paulus/chirpy/internal/auth"
	"github.com/jan-paulus/chirpy/internal/database"
)

func (cfg *apiConfig) webhookPolka(w http.ResponseWriter, r *http.Request) {
	type response struct{}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	providedApiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		if errors.Is(err, auth.ErrNoAuthHeaderIncluded) {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve auth header")
		return
	}

	if providedApiKey != cfg.polkaApiKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	switch params.Event {
	case "user.upgraded":
		_, err := cfg.DB.UpradeUser(params.Data.UserId)
		if err != nil {
			if errors.Is(err, database.ErrNotExist) {
				respondWithError(w, http.StatusNotFound, "Couldn't find user")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
			return
		}
	}

	respondWithJSON(w, http.StatusOK, response{})
}
