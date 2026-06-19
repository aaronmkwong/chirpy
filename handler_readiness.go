// healthzHandler responds to readiness checks
// returns 200 OK with a plain text body to indicate the server is ready

package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	
	// set a response header
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// write a status code
	w.WriteHeader(http.StatusOK)

	// write a body
	w.Write([]byte("OK"))
}
