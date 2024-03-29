package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tcluri/chirpy/internal/database"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
	polkaSecret    string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Load the environment variable
	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable not set")
	}

	// Welcome message
	fmt.Println("Hello! Welcome to the chirpy webserver!")

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}

	// Initialize the api config
	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:      jwtSecret,
		polkaSecret:    polkaKey,
	}
	// mux := http.NewServeMux()

	router := chi.NewRouter() // app router
	// fsHandler := apiCfg.middlewareMetricsInc(middlewareLog(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	router.Mount("/", apiCfg.middlewareMetricsInc(middlewareLog(http.FileServer(http.Dir(".")))))

	router.Get("/metrics", apiCfg.handlerMetrics)

	// API router endpoints
	apiRouter := chi.NewRouter() // api router
	apiRouter.Get("/healthz", handlerReadiness)

	apiRouter.Post("/chirps", apiCfg.handlerChirpsCreate)
	apiRouter.Get("/chirps", apiCfg.handlerChirpsRetrieve)
	apiRouter.Get("/chirps/{chirpID}", apiCfg.handlerChirpsGet)
	apiRouter.Delete("/chirps/{chirpID}", apiCfg.handlerChirpsDelete)

	apiRouter.Put("/users", apiCfg.handlerUsersUpdate)
	apiRouter.Post("/users", apiCfg.handlerUsersCreate)

	apiRouter.Post("/polka/webhooks", apiCfg.handlerUserUpgrade)

	apiRouter.Post("/refresh", apiCfg.handlerRefresh)
	apiRouter.Post("/revoke", apiCfg.handlerRevoke)
	apiRouter.Post("/login", apiCfg.handlerUsersLogin)
	router.Mount("/api", middlewareLog(apiRouter))

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}
	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
