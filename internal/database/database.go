package database

import (
	"encoding/json"
	"errors"
	"os"
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

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Hash  []byte `json:"hash"`
}

var ErrAlreadyExists = errors.New("User already exists")

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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
		Users:  make(map[int]User),
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

func (db *DB) ResetDB() error {
	path := db.path
	if err := os.Remove(path); err != nil {
		return err
	}
	newDB, err := NewDB(path)
	if err != nil {
		return err
	}
	db = newDB
	return nil
}
