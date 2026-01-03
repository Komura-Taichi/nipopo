package server

import (
	"log"
	"net/http"

	"github.com/Komura-Taichi/nipopo/backend/internal/handler"
)

func main() {
	http.HandleFunc("/healthz", handler.Healthz)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("server start error: %v", err)
	}
}
