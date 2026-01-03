package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Komura-Taichi/nipopo/backend/internal/handler"
)

type healthzResponse struct {
	Status string `json:"status"`
}

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/healthz", nil)

	handler.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, body=%s", rec.Code, rec.Body.String())
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type=%q", ct)
	}

	var res healthzResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("invalid json=%v, body=%s", err, rec.Body.String())
	}

	if res.Status != "OK" {
		t.Fatalf("status field: want=%q, got=%q", "OK", res.Status)
	}
}
