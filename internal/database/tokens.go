package database

import (
	"errors"
	"time"
)

func (db *DB) RevokeToken(tokenToRevoke string) error {
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	alreadyRevoked, err := db.IsTokenRevoked(tokenToRevoke)
	if alreadyRevoked {
		return errors.New("Token already revoked")
	}
	// Create a new revoked token and write to db
	revoked := RevokedToken{
		ID:        tokenToRevoke,
		RevokedAt: time.Now().UTC(),
	}
	dbStruct.RevokedTokens[tokenToRevoke] = revoked
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) IsTokenRevoked(tokenToCheck string) (bool, error) {
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return false, err
	}
	_, ok := dbStruct.RevokedTokens[tokenToCheck]
	if ok {
		return true, nil
	}
	return false, nil
}
