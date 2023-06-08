package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	// Initialize the mutex
	mux := &sync.RWMutex{}
	// Create an empty DBStructure
	dbStruct := DBStructure{}

	// Read the database file
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Create an empty database file if it doesn't exist
			data, err = createEmptyDatabaseFile(path)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	// Unmarshall the data - with data or otherwise
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return nil, err
	}

	// Return the new DB instance
	return &DB{
		path: path,
		mux:  mux,
	}, nil
}

func createEmptyDatabaseFile(path string) ([]byte, error) {
	emptyDB := DBStructure{
		Chirps: make(map[int]Chirp),
	}

	data, err := json.MarshalIndent(emptyDB, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func (db *DB) ensureDB() error {
	_, err := os.Stat(db.path)
	if errors.Is(err, os.ErrNotExist) {
		_, err = createEmptyDatabaseFile(db.path)
		return err
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStruct := DBStructure{}
	err = json.Unmarshal(data, &dbStruct)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Generate a unique ID for the chirp
	id := db.generateUniqueID()

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

// generateUniqueID generates a unique ID for the chirp
func (db *DB) generateUniqueID() int {
	dbStruct, _ := db.loadDB()
	maxID := 0
	for id := range dbStruct.Chirps {
		if id > maxID {
			maxID = id
		}
	}
	return maxID + 1
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
