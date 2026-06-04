package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/erikw/gomemo/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Starting Gomemo with\n%v\n", cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(20 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		// fmt.Println(middleware.GetReqID(r.Context()))
		// _, err := w.Write([]byte("Hello, world"))
		// if err != nil {
		//	// TODO use structured logging pkg
		//	fmt.Printf("Error serving request: %v\n", err)
		// }
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})
	err := http.ListenAndServe(cfg.AddrString(), r)
	if err != nil {
		fmt.Printf("Error serving: %v\n", err)
	}
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
