package notes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/erikw/gomemo/internal/httpx"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	logger  *slog.Logger
	service *Service
}

type createNoteRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func NewHandler(logger *slog.Logger, service *Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/notes", h.handleGetAll) // TODO pagingate
	r.Post("/notes", h.handleCreate)
	r.Get("/notes/{noteID}", h.handleGetByID)
	r.Patch("/notes/{noteID}", h.handleUpdateByID)
	r.Delete("/notes/{noteID}", h.handleDeleteByID)
}

func (h *Handler) handleGetAll(w http.ResponseWriter, req *http.Request) {

	notes, err := h.service.GetAll(req.Context())
	if err != nil {
		h.logger.Error("could not fetch Notes")
		httpx.RespondError(w, http.StatusNotFound, "Note could not be found.")
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, notes); err != nil {
		h.logger.Error("could not respond with JSON encoding", "notes", notes)
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, req *http.Request) {
	noteReq, err := h.noteFromRequest(w, req)
	if err != nil {
		return
	}

	note, err := h.service.Create(req.Context(), *noteReq.Title, *noteReq.Content)
	if err != nil {
		h.logger.Error("could not create Note", "error", err.Error())
		switch {
		case errors.Is(err, ErrTitleRequired):
			httpx.RespondError(w, http.StatusBadRequest, err.Error())
		default:
			httpx.RespondError(w, http.StatusInternalServerError, "Error creating Note.")
		}
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, note); err != nil {
		h.logger.Error("could not respond with JSON encoding", "noteID", note.ID, "note", note)
	}

}

func (h *Handler) handleGetByID(w http.ResponseWriter, req *http.Request) {
	id, err := h.idFromPath(w, req)
	if err != nil {
		return
	}

	note, err := h.service.GetByID(req.Context(), id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not fetch Note with ID `%d`", id))
		httpx.RespondError(w, http.StatusNotFound, "Note could not be found.")
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, note); err != nil {
		h.logger.Error("could not respond with JSON encoding", "noteID", id, "note", note)
	}
}

func (h *Handler) handleUpdateByID(w http.ResponseWriter, req *http.Request) {
	id, err := h.idFromPath(w, req)
	if err != nil {
		return
	}

	noteReq, err := h.noteFromRequest(w, req)
	if err != nil {
		return
	}

	note, err := h.service.Update(req.Context(), id, noteReq.Title, noteReq.Content)
	if err != nil {
		errReq := fmt.Errorf("could not parse save updated Note")
		h.logger.Error(errReq.Error())
		httpx.RespondError(w, http.StatusBadRequest, errReq.Error())
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, note); err != nil {
		h.logger.Error("could not respond with JSON encoding", "noteID", id, "note", note)
	}
}

func (h *Handler) handleDeleteByID(w http.ResponseWriter, req *http.Request) {
	id, err := h.idFromPath(w, req)
	if err != nil {
		return
	}

	deleted, err := h.service.DeleteByID(req.Context(), id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not delete Note with ID `%d`", id))
		httpx.RespondError(w, http.StatusNotFound, "Note could not be deleted.")
		return
	}

	if deleted {
		if err = httpx.RespondJSON(w, http.StatusNoContent, nil); err != nil {
			h.logger.Error("could not respond with JSON encoding", "noteID", id)
		}
	} else {
		if err = httpx.RespondJSON(w, http.StatusNotFound, nil); err != nil {
			h.logger.Error("could not respond with JSON encoding", "noteID", id)
		}
	}

}

func (h *Handler) idFromPath(w http.ResponseWriter, req *http.Request) (int64, error) {
	idStr := chi.URLParam(req, "noteID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not convert urlParam `%s` noteID to int64", idStr))
		httpx.RespondError(w, http.StatusBadRequest, "Invalid notesID.")
		return -1, fmt.Errorf("could not extract Note ID from URL query parameters")
	}
	return id, nil
}

func (h *Handler) noteFromRequest(w http.ResponseWriter, req *http.Request) (createNoteRequest, error) {
	var noteReq createNoteRequest

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&noteReq); err != nil {
		if errors.Is(err, io.EOF) {
			return noteReq, nil
		}
		errReq := fmt.Errorf("could not parse Note JSON from request body: %s", err.Error())
		h.logger.Error(errReq.Error())
		httpx.RespondError(w, http.StatusBadRequest, "Invalid JSON format for Note.")
		return createNoteRequest{}, errReq
	}
	return noteReq, nil
}
