
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
)

// hold any stateful, in-memory data
type apiConfig struct {
	fileserverHits atomic.Int32 // safely increment, read integer value across multiple goroutines (HTTP requests)
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

func main() {

	// instantiate api config
	apiCfg := apiConfig{
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

	// define configuration and behavior for running an active HTTP server
	serveStruct := http.Server{
		Addr: ":8080",
		Handler: serveMux,
	}

	// start the server
	err := serveStruct.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

