package handler

import (
	"net/http"

	"github.com/Komura-Taichi/nipopo/backend/internal/handler/middleware"
)

func requireUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || userID == "" {
		writeErrorJSON(w, http.StatusUnauthorized, "unauthorized", nil)
		return "", false
	}

	return userID, true
}
