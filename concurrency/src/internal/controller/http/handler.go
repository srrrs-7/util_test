package http

import (
	"concurrency/internal/controller/router"
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
		r.Post("/create", rt.concurrent.Create)
		// check
		r.Get("/check/{id}", rt.concurrent.Check)
	})

	return router
}
