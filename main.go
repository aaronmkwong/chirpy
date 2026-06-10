package main

import (
	"net/http"
	"fmt"
	"os"
)

func main() {
	
	// create HTTP request router and multiplexer 
	newServeMux := http.NewServeMux()

	// define configuration and behavior for running an active HTTP server
	serveStruct := http.Server{
		Addr: ":8080",
		Handler: newServeMux,
	}

	// start the server
	err := serveStruct.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

