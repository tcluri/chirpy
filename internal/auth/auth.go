package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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

func CreateJWT(userid int, jwtSecret string, expirytime time.Duration) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expirytime)),
		Subject:   string(userid),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GetBearerToken(header http.Header) (string, error) {
	authField := header.Get("Authorization")
	fmt.Printf("The authorization token is %v", authField)
	token := strings.TrimPrefix(authField, "Bearer: ")
	fmt.Printf("The token extracted %v", token)
	if token == "" {
		return "", errors.New("Couldn't find Token in the request header")
	}
	return token, nil
}

func ValidateJWT(tokenString string, jwtSecret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		if claims.ExpiresAt.Before(time.Now().UTC()) {
			return "", errors.New("Token expired")
		}
		userID := claims.Subject
		return userID, nil
	} else {
		return "", errors.New("Couldn't validate JWT token")
	}
}
