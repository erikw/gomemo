package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/erikw/gomemo/internal/api"
	"github.com/erikw/gomemo/internal/config"
	"github.com/erikw/gomemo/internal/notes"
	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/testutil"
)

func TestNotesLifecycleIntegration(t *testing.T) {
	logger := testutil.Logger()
	router := api.NewRouter(logger, config.Config{Host: "127.0.0.1", Port: "8080"})
	store := storage.NewMemory[*notes.Note](logger)
	service := notes.NewService(logger, store)
	handler := notes.NewHandler(logger, service)
	router.RegisterV1Routes(handler)

	createBody := bytes.NewBufferString(`{"title":"integration note","content":"hello"}`)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/notes", createBody)
	createRec := httptest.NewRecorder()
	router.ChiRouter().ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("POST /api/v1/notes status = %d, want %d", createRec.Code, http.StatusOK)
	}

	var created notes.Note
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response error = %v", err)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/notes?limit=10&offset=0", nil)
	listRec := httptest.NewRecorder()
	router.ChiRouter().ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/v1/notes status = %d, want %d", listRec.Code, http.StatusOK)
	}

	updateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/notes/"+strconv.FormatInt(created.ID, 10), bytes.NewBufferString(`{"title":"updated"}`))
	updateRec := httptest.NewRecorder()
	router.ChiRouter().ServeHTTP(updateRec, updateReq)
	if updateRec.Code != http.StatusOK {
		t.Fatalf("PATCH /api/v1/notes/{id} status = %d, want %d", updateRec.Code, http.StatusOK)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/notes/"+strconv.FormatInt(created.ID, 10), nil)
	deleteRec := httptest.NewRecorder()
	router.ChiRouter().ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /api/v1/notes/{id} status = %d, want %d", deleteRec.Code, http.StatusNoContent)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/notes/"+strconv.FormatInt(created.ID, 10), nil)
	getRec := httptest.NewRecorder()
	router.ChiRouter().ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusNotFound {
		t.Fatalf("GET /api/v1/notes/{id} after delete status = %d, want %d", getRec.Code, http.StatusNotFound)
	}
}
