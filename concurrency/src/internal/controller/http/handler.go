package http

import (
	"concurrency/internal/controller/server"
	"net/http"

	"github.com/go-chi/chi"
)

func Serve() *chi.Mux {
	router := chi.NewMux()

	router.Use(router.Middlewares()...)

	// health
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {})

	// auth domain path
	router.Route("/concurrent/v1", func(r chi.Router) {
		// create
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			c := server.Create{}
			c.Create(w, r.WithContext(r.Context()))
		})

		// check
		r.Get("/check/{id}", func(w http.ResponseWriter, r *http.Request) {
			c := server.Check{}
			c.Check(w, r.WithContext(r.Context()))
		})
	})

	return router
}
