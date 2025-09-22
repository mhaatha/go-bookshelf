package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

// ReadFromRequestBody reads the request body and stores it in the result parameter
func ReadFromRequestBody(r *http.Request, result interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}

	return nil
}

// WriteToResponseBody writes the result parameter to the response body
func WriteToResponseBody(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
