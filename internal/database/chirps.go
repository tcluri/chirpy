package database

import (
	"errors"
	"sort"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	// Generate a unique ID for the chirp
	id := len(dbStruct.Chirps) + 1
	// Create the chirp
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	// Add the chirp to the database
	dbStruct.Chirps[id] = chirp
	// Write the updated database back to disk
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	// Prepare the chirps slice
	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	// Append chirps to the slice
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}
	// Sort chirps by ID in ascending order
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	return chirps, nil
}

func (db *DB) GetChirp(chirpID int) (Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp := dbStruct.Chirps[chirpID]
	if chirp.ID == 0 && chirp.Body == "" {
		return Chirp{}, errors.New("The chirp does not exist")
	}
	return chirp, nil
}
