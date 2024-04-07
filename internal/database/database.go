package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return db, err
	}

	log.Printf("Created new DB: %s", db.path)
	return db, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%s does not exist, creating new...", db.path)

		emptyDb := DBStructure{
			Chirps:        map[int]Chirp{},
			Users:         map[int]User{},
			RefreshTokens: map[string]RefreshToken{},
		}
		contents, err := json.Marshal(emptyDb)

		if err != nil {
			return err
		}

		err = os.WriteFile(db.path, contents, 0666)
		if err != nil {
			return errors.New("Failed to create database file.")
		}
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	contents, err := os.ReadFile(db.path)

	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(contents, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	contents, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, contents, 0666)
	return err
}
