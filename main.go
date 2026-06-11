
// in CLI compile to temporary binary and execute (when developing)
// go run main.go

// in CLI rebuild and run 
// go build -o out && ./out


// in CLI send signal to server to stop and exit
// CTRL + C

package main

import (
	"net/http"
	"log"
)

// healthzHandler responds to readiness checks.
// Returns 200 OK with a plain text body to indicate the server is ready.
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	
	// set a response header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// write a status code
	w.WriteHeader(http.StatusOK)

	// write a body
	w.Write([]byte("OK"))
}

func main() {

	// create directory reference
	dirRef := http.Dir(".")
	fileServHandler:= http.FileServer(dirRef)
	
	// create HTTP request router and multiplexer 
	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", fileServHandler))

	// register healthz handler
	serveMux.HandleFunc("/healthz", healthzHandler)

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

