
// in CLI compile to temporary binary and execute (when developing)
// go run main.go

// in CLI rebuild and run 
// go build -o out && ./out

package main

import (
	"net/http"
	"log"
)

func main() {
	
	// create directory reference
	dirRef := http.Dir(".")
	fileServHandler:= http.FileServer(dirRef)
	
	// create HTTP request router and multiplexer 
	serveMux := http.NewServeMux()
	serveMux.Handle("/", fileServHandler)

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

