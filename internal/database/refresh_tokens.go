package database

import (
	"time"
)

type RefreshToken struct {
	Token     string `json:"token"`
	RevokedAt string `json:"revoked_at"`
	Revoked   bool   `json:"revoked"`
}

func (db *DB) CreateRefreshToken(token string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken := RefreshToken{
		Token:   token,
		Revoked: false,
	}

	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

func (db *DB) RevokeRefreshToken(token string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return RefreshToken{}, ErrNotExist
	}

	refreshToken.RevokedAt = time.Now().UTC().Format(time.UnixDate)
	refreshToken.Revoked = true

	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return RefreshToken{}, err
	}

	return refreshToken, nil
}

func (db *DB) GetRefreshToken(token string) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if ok {
		return refreshToken, nil
	}

	return RefreshToken{}, ErrNotExist
} 
