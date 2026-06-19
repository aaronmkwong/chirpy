// reset method

package main

import "net/http" 

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// check platform/environment
	if cfg.platform != "dev" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset only allowed in development environment"))
		return
	}

	// delete all users from the database using SQLC method
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete users"))
		return
	}

	// reset fileserverHits to 0 before writing the response
	cfg.fileserverHits.Store(0)

	// set the success response header and status code
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset and database cleared"))
}