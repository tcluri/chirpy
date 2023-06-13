package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tcluri/chirpy/internal/database"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
}

func main() {
	fmt.Println("Hello! Welcome to the chirpy webserver!")

	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println("We are here")
		log.Fatal(err)
	}
	// Initialize the api config
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
	}
	// mux := http.NewServeMux()

	router := chi.NewRouter() // app router
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)
	// API router endpoints
	apiRouter := chi.NewRouter() // api router
	// apiRouter.Post("/users", apiCfg.handlerUsersCreate)
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	router.Mount("/api", apiRouter)
	// Admin router endpoints
	adminRouter := chi.NewRouter()
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)
	router.Mount("/admin", adminRouter)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
