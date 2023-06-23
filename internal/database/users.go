package database

import "errors"

func (db *DB) CreateUser(email string, hashedPassword []byte) (User, error) {
	// Load the current database
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	// Check if the user already exists
	existingUser, err := db.GetUserByEmail(email)
	if err == nil && len(existingUser.Email) > 0 {
		return User{}, ErrAlreadyExists
	}
	// Generate a unique ID for the user
	id := len(dbStruct.Users) + 1
	// IsChirpyRed subscribed
	subscribed := false
	// Create the user
	user := User{
		ID:          id,
		Email:       email,
		Hash:        hashedPassword,
		IsChirpyRed: subscribed,
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
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStruct.Users[userIDInt]
	if !ok {
		return User{}, errors.New("User does not exist")
	}
	user.Email = email
	user.Hash = hashedPassword
	dbStruct.Users[userIDInt] = user
	// Write the changes to disk
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeUserStatus(userIDInt int) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	user, ok := dbStruct.Users[userIDInt]
	if !ok {
		return User{}, errors.New("User does not exist")
	}
	user.IsChirpyRed = true
	dbStruct.Users[userIDInt] = user
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
