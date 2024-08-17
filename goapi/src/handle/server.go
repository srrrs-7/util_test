package handle

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	HEALTH_PATH = "/health"
	DOMAIN_PATH = "/domain/v1"
)

func NewServer() *chi.Mux {
	r := chi.NewMux()

	r.Use(r.Middlewares()...)

	r.Get(HEALTH_PATH, func(w http.ResponseWriter, r *http.Request) {})
	r.Get(DOMAIN_PATH, func(w http.ResponseWriter, r *http.Request) {})

	return r
}
