package main

import (
	"fmt"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Hello! Welcome to the chirpy webserver!")

	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)

	newServer := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}
	newServer.ListenAndServe()
}
