package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jan-paulus/chirpy/internal/auth"
)

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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	defaultExpiration := 60 * 60 * 24
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.Id, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Id:    user.Id,
		Email: user.Email,
		Token: token,
	})

	// expiresAt := &jwt.NumericDate{}
	//
	// if params.ExpiresInSeconds != 0 {
	// 	expiresAt = jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(params.ExpiresInSeconds)))
	// } else {
	// 	expiresAt = jwt.NewNumericDate(time.Now().UTC().Add(time.Hour * 24))
	// }
	//
	// claims := &jwt.RegisteredClaims{
	// 	Issuer:    "chirpy",
	// 	Subject:   strconv.Itoa(user.Id),
	// 	IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
	// 	ExpiresAt: expiresAt,
	// }
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// signedString, err := token.SignedString([]byte(cfg.jwtSecret))
	//
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "Failed to authorize.")
	// 	return
	// }
	//
	// responseBody := response{
	// 	Email: user.Email,
	// 	Id:    user.Id,
	// 	Token: signedString,
	// }
	// respondWithJSON(w, http.StatusOK, responseBody)
}
