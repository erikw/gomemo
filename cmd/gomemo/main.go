package main

import (
	"fmt"
	"net/http"

	"github.com/erikw/gomemo/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	fmt.Printf("Starting Gomemo with\n%v\n", cfg)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello, world"))
		if err != nil {
			// TODO use structured logging pkg
			fmt.Printf("Error serving request: %v\n", err)
		}
	})
	err := http.ListenAndServe(cfg.Port, r)
	if err != nil {
		fmt.Printf("Error serving: %v\n", err)
	}
}
