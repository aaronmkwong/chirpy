// in CLI compile to temporary binary and execute (when developing)
// go run main.go

// in CLI rebuild and run
// go build -o out && ./out

// in CLI send signal to server to stop and exit
// CTRL + C

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/aaronmkwong/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver anonymously for its side effects (registering the driver)
)

// hold any stateful, in-memory data
type apiConfig struct {
	fileserverHits atomic.Int32 // safely increment, read integer value across multiple goroutines (HTTP requests)
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func main() {

	// load environment variables safely
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// verify the connection works
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// initialize queries
	dbQueries := database.New(db)

	// load the JWT secret from .env file
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	// load platform from .env file
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	// instantiate api config
	apiCfg := apiConfig{
		db:             dbQueries,
		fileserverHits: atomic.Int32{},
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	// create directory reference
	dirRef := http.Dir(".")
	fileServHandler := http.FileServer(dirRef)

	// create HTTP request router and multiplexer
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServHandler)))

	// register healthz handler
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	// register request handler
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// register reset handler
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	// register chirps validate handler
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	// register create user handler
	serveMux.HandleFunc("POST /api/users", apiCfg.handlerUserCreate)

	// register get many chirps handler
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)

	// register get single chirp handler
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGet)

	// register login handler
	serveMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	// register refresh handler
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	// register revoke handler
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	// register update user handler
	serveMux.HandleFunc("PUT /api/users", apiCfg.handlerUserUpdate)

	// define configuration and behavior for running an active HTTP server
	serveStruct := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	// start the server
	err = serveStruct.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
