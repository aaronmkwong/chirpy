// reset method

package main

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
    // set the response header and status code
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

    // reset fileserverHits to 0 using the atomic method:
    cfg.fileserverHits.Store(0)
}
