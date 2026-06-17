package notes

import (
	"fmt"
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

func NewHandler(logger *slog.Logger, service *Service) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/notes", h.HandleGetAll) // TODO pagingate
	r.Get("/notes/{noteID}", h.HandleGetByID)
	// TODO
	// POST   /notes
	// DELETE /notes/{id}
}
func (h *Handler) HandleGetAll(w http.ResponseWriter, req *http.Request) {

	notes, err := h.service.GetAll(req.Context())
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not fetch Notes"))
		httpx.RespondError(w, http.StatusNotFound, "Note could not be found")
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, notes); err != nil {
		h.logger.Error("could not respond with JSON encoding", "notes", notes)
	}
}

func (h *Handler) HandleGetByID(w http.ResponseWriter, req *http.Request) {
	idStr := chi.URLParam(req, "noteID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not convert urlParam `%s` noteID to int64", idStr))
		httpx.RespondError(w, http.StatusBadRequest, "Invalid notesID")
		return
	}

	note, err := h.service.GetByID(req.Context(), id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("could not fetch Note with ID `%d`", id))
		httpx.RespondError(w, http.StatusNotFound, "Note could not be found")
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, note); err != nil {
		h.logger.Error("could not respond with JSON encoding", "noteID", id, "note", note)
	}
}
