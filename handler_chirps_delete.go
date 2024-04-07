package main

import (
	"net/http"
	"strconv"

	"github.com/jan-paulus/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	type response struct{}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	issuer, subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	if issuer != "chirpy-access" {
		respondWithError(w, http.StatusUnauthorized, "Provided token is not a refresh token")
		return
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

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
		respondWithError(w, http.StatusInternalServerError, "Failed to find chirp")
		return
	}

	if chirp.AuthorId != userId {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	err = cfg.DB.DeleteChirp(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, response{})
}
