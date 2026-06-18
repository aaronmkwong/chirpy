// metrics method 
// writes the number of requests that have been counted as plain text in this format to the HTTP response

package main

import (
    "fmt"
    "net/http"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	// set the response header and status code
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

    // writes the number of requests that have been counted as plain text in this format to the HTTP response
	hits := cfg.fileserverHits.Load()
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", hits)))
}

// middleware method that increments the fileserverHits counter every time called
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}

