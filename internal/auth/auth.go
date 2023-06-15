package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	// Convert password string into bytes and then hash it
	passBytes := []byte(password)
	hashedPass, err := bcrypt.GenerateFromPassword(passBytes, 10)
	if err != nil {
		return []byte{}, err
	}
	return hashedPass, nil
}

func CheckPasswordHash(password string, hashedPass []byte) error {
	// Check if the hash and password match
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
	if err != nil {
		return err
	}
	return nil
}
