package httpx

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	err := RespondJSON(rec, http.StatusCreated, map[string]string{"status": "ok"})
	if err != nil {
		t.Fatalf("RespondJSON() error = %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body error = %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body status = %q, want %q", body["status"], "ok")
	}
}

func TestRespondError(t *testing.T) {
	rec := httptest.NewRecorder()
	RespondError(rec, http.StatusBadRequest, "bad request")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body error = %v", err)
	}
	if body["message"] != "bad request" {
		t.Fatalf("body message = %q, want %q", body["message"], "bad request")
	}
}
