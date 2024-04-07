package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/jan-paulus/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

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

	if issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "Provided token is not a refresh token")
		return
	}

	userId, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	refreshToken, err := cfg.DB.GetRefreshToken(token)
	if refreshToken.Revoked {
		respondWithError(w, http.StatusUnauthorized, "Refresh token was revoked")
		return
	}

	accessToken, err := auth.MakeJWT("chirpy-access", userId, cfg.jwtSecret, time.Duration(1)*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}
