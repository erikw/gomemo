package notes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/erikw/gomemo/internal/storage"
	"github.com/erikw/gomemo/internal/testutil"
	"github.com/go-chi/chi/v5"
)

func newNotesRouterForTest() chi.Router {
	logger := testutil.Logger()
	store := storage.NewMemory[*Note](logger)
	service := NewService(logger, store)
	handler := NewHandler(logger, service)
	r := chi.NewRouter()
	handler.RegisterRoutes(r)
	return r
}

func TestHandleCreateAndGetByID(t *testing.T) {
	r := newNotesRouterForTest()

	createReq := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(`{"title":"my note","content":"hello"}`))
	createRec := httptest.NewRecorder()
	r.ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusOK {
		t.Fatalf("POST /notes status = %d, want %d", createRec.Code, http.StatusOK)
	}

	var created Note
	if err := json.NewDecoder(createRec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response error = %v", err)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/notes/"+strconv.FormatInt(created.ID, 10), nil)
	getRec := httptest.NewRecorder()
	r.ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET /notes/{id} status = %d, want %d", getRec.Code, http.StatusOK)
	}

	var got Note
	if err := json.NewDecoder(getRec.Body).Decode(&got); err != nil {
		t.Fatalf("decode get response error = %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("GET /notes/{id} ID = %d, want %d", got.ID, created.ID)
	}
}

func TestHandleCreateInvalidJSON(t *testing.T) {
	r := newNotesRouterForTest()

	req := httptest.NewRequest(http.MethodPost, "/notes", strings.NewReader(`{"title":"broken"`))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("POST /notes invalid JSON status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleGetAllInvalidPagination(t *testing.T) {
	r := newNotesRouterForTest()

	req := httptest.NewRequest(http.MethodGet, "/notes?limit=abc", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("GET /notes invalid pagination status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
