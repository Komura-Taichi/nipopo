package handler

import (
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{ "status": "OK" }`)); err != nil {
		log.Printf("healthz write error: %v", err)
	}
}
