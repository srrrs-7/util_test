package http

import (
	"concurrency/internal/controller/router"
	"context"
	"net/http"

	"github.com/go-chi/chi"
)

type Router struct {
	concurrent *router.Concurrent
}

func (rt *Router) Serve() *chi.Mux {
	router := chi.NewMux()

	router.Use(router.Middlewares()...)

	// health
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {})

	router.Route("/concurrent/v1", func(r chi.Router) {
		// create
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			rt.concurrent.Create(w, r)
		})
		// check
		r.Get("/check/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			r = r.WithContext(context.WithValue(r.Context(), "id", id))
			rt.concurrent.Check(w, r)
		})
	})

	return router
}
