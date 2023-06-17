package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/tcluri/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password        string `json:"password"`
		Email           string `json:"email"`
		ExpiryInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	fmt.Printf("Parameters are %v and expiry in %v and password is %v\n", params.Email, params.ExpiryInSeconds, params.Password)
	user, err := cfg.DB.GetUserByEmail(params.Email)
	fmt.Printf("Handler login: User is %v\n", user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.Hash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	fmt.Println("Password check complete!")

	defaultExpiration := 60 * 60 * 24
	if params.ExpiryInSeconds == 0 {
		params.ExpiryInSeconds = defaultExpiration
	} else if params.ExpiryInSeconds > defaultExpiration {
		params.ExpiryInSeconds = defaultExpiration
	}

	fmt.Printf("Login Handler: The expiration time now is: %v\n", params.ExpiryInSeconds)
	token, err := auth.CreateJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiryInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}
	fmt.Printf("The token is: %v\n", token)
	fmt.Println("Login Handler: Token is got and we are sending the success response")
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}