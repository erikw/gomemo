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
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) HandleGetByID(w http.ResponseWriter, req *http.Request) {
	idStr := chi.URLParam(req, "noteID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not convert urlParam `%s` noteID to int64.", idStr))
		httpx.RespondError(w, http.StatusBadRequest, "Invalid notesID")
		return
	}

	note, err := GetByID(id)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Could not fetch Note with ID `%d`.", id))
		httpx.RespondError(w, http.StatusNotFound, "Note could not be found")
		return
	}

	if err = httpx.RespondJSON(w, http.StatusOK, note); err != nil {
		h.logger.Error("Could not respond with JSON encoding", "noteID", id, "note", note)
	}
}
