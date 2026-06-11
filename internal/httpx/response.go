package httpx

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	resp := map[string]string{
		"message": message,
	}
	_ = RespondJSON(w, status, resp)
}
