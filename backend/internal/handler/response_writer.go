package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeErrorJSON(w http.ResponseWriter, status int, msg string, details map[string]any) {
	response := ErrorResponse{
		Error: ErrorBody{
			Code:    status,
			Message: msg,
			Details: details,
		},
	}

	writeJSON(w, status, response)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(v); err != nil {
		log.Printf("error happened while encoding JSON: %v", err)
	}
}
