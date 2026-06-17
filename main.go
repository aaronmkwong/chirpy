
// in CLI compile to temporary binary and execute (when developing)
// go run main.go

// in CLI rebuild and run 
// go build -o out && ./out


// in CLI send signal to server to stop and exit
// CTRL + C

package main

import (
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
	"encoding/json"
	"strings"
	"database/sql"
	"os"
	"github.com/joho/godotenv"
	"github.com/aaronmkwong/chirpy/internal/database"
	_ "github.com/lib/pq"  // Import the PostgreSQL driver anonymously for its side effects (registering the driver)
)

// hold any stateful, in-memory data
type apiConfig struct {
	fileserverHits atomic.Int32 // safely increment, read integer value across multiple goroutines (HTTP requests)
	DB *database.Queries
}

// middleware method that increments the fileserverHits counter every time called
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

// reset method
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
    // set the response header and status code
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

    // reset fileserverHits to 0 using the atomic method:
    cfg.fileserverHits.Store(0)
}

// metrics method 
// writes the number of requests that have been counted as plain text in this format to the HTTP response
func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	// set the response header and status code
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

    // writes the number of requests that have been counted as plain text in this format to the HTTP response
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)))
}

// healthzHandler responds to readiness checks
// returns 200 OK with a plain text body to indicate the server is ready
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	
	// set a response header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// write a status code
	w.WriteHeader(http.StatusOK)

	// write a body
	w.Write([]byte("OK"))
}

// replaces forbidden words with asterisks
// used by handlerChirpsValidate
func CleanChirp(chirp string) string {
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(chirp, " ")

	for i, word := range words {
		if badWords[strings.ToLower(word)] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

// decodes and validates incoming chirp body and ensures chirp is 140 characters or less
func handlerChirpsValidate(w http.ResponseWriter, r *http.Request){

	// parsing the incoming request
	type parameters struct {
    	Body string `json:"body"`
	}


	// sending back an error response
	type errorResponse struct {
		Error string `json:"error"`
	}	

	// cleaning bad words and success response 
	type successResponse struct {
    	CleanedBody string `json:"cleaned_body"`  
	}	
	
    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	if len(params.Body) > 140 {
		// create error response struct
		respBody := errorResponse{
			Error: "Chirp is too long",
		}
		
		// marshal to JSON bytes
		dat, err := json.Marshal(respBody)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		
		// write headers and the response data back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400) // 400 Bad Request
		w.Write(dat)
		return
	}
	
	// clean profane words and build success response
		respBody := successResponse{
		CleanedBody: CleanChirp(params.Body),
	}

    // marshal to JSON bytes
    dat, err := json.Marshal(respBody)
    if err != nil {
        w.WriteHeader(500)
        return
    }

    // write headers and response data back to the client
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(200)
    w.Write(dat)	
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

	// instantiate api config
	apiCfg := apiConfig{
		DB: dbQueries,
		fileserverHits: atomic.Int32{},
	}

	// create directory reference
	dirRef := http.Dir(".")
	fileServHandler:= http.FileServer(dirRef)
	
	// create HTTP request router and multiplexer 
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServHandler)))

	// register healthz handler
	serveMux.HandleFunc("GET /api/healthz", healthzHandler)

	// register request handler
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)

	// register reset handler
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	// register chirps validate handler 
	serveMux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	// define configuration and behavior for running an active HTTP server
	serveStruct := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	// start the server
	err = serveStruct.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

