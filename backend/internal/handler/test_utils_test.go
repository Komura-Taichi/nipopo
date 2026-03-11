package handler_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("status=%d want=%d", rec.Code, want)
	}
}

func assertContentType(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("content-type=%q", ct)
	}
}

func unmarshalJSON(t *testing.T, data []byte, v any) {
	t.Helper()
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("invalid json: %v body=%s", err, string(data))
	}
}
