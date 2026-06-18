// decodes and validates incoming chirp body and ensures chirp is 140 characters or less

package main

import (
	"net/http"
	"log"
	"encoding/json"
	"strings"
)

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
		CleanedBody: getCleanedChirp(params.Body),
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

// replaces forbidden words with asterisks
// used by handlerChirpsValidate
func getCleanedChirp(chirp string) string {
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