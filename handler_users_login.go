package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/tcluri/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.Hash)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Access token chirpy-access
	access_issuer := "chirpy-access"
	access_expiry := 60 * 60
	access_token, err := auth.CreateJWT(user.ID, cfg.jwtSecret, time.Duration(access_expiry)*time.Second, access_issuer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT for access")
		return
	}

	// Refresh token chirpy-refresh
	refresh_issuer := "chirpy-refresh"
	refresh_expiry := 60 * 60 * 24 * 60
	refresh_token, err := auth.CreateJWT(user.ID, cfg.jwtSecret, time.Duration(refresh_expiry)*time.Second, refresh_issuer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT for refresh")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
		Token:        access_token,
		RefreshToken: refresh_token,
	})
}
