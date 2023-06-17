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

func (db *DB) CreateUser(email string, hashedPassword []byte) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Check if the user already exists
	existingUser, err := db.GetUserByEmail(email)
	if len(existingUser.Email) > 0 {

		return User{}, ErrAlreadyExists
	}

	// Generate a unique ID for the user
	id := len(dbStruct.Users) + 1

	// Create the user
	user := User{
		ID:    id,
		Email: email,
		Hash:  hashedPassword,
	}

	// Add the user to the database
	dbStruct.Users[id] = user

	// Write the updated database back to disk
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(useremail string) (User, error) {
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	// Find the user
	for _, eachUser := range dbStruct.Users {
		if eachUser.Email == useremail {
			return eachUser, nil
		}
	}
	return User{}, errors.New("Could not find user")
}

func (db *DB) UpdateUser(userIDInt int, email string, hashedPassword []byte) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, eachUser := range dbStruct.Users {
		if eachUser.ID == userIDInt {
			// Update the user in db
			dbStruct.Users[userIDInt] = User{
				ID:    userIDInt,
				Email: email,
				Hash:  hashedPassword,
			}
			return dbStruct.Users[userIDInt], nil
		}
	}
	return User{}, errors.New("User not available in the database")
}
