package main

import (
	"encoding/json"
	"net/http"

	"github.com/tcluri/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetPolkaApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}
	if apiKey != cfg.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "API key mismatch")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	// Status OK if it is any other event
	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusOK, struct{}{})
		return
	} else {
		_, err = cfg.DB.UpgradeUserStatus(params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Could not be found")
			return
		}
		respondWithJSON(w, http.StatusOK, struct{}{})
	}
}
