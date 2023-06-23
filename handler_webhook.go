package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUserUpgrade(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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
