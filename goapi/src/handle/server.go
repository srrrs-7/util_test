package handle

import (
	"api/domain/usecase"
	"api/util/static"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	checker usecase.CheckUseCase
	creator usecase.CreateUseCase
}

func NewServer(checker usecase.CheckUseCase, creator usecase.CreateUseCase) Server {
	return Server{checker, creator}
}

func (s Server) Routing() *chi.Mux {
	r := chi.NewMux()

	r.Use(r.Middlewares()...)

	r.Get(static.HEALTH_PATH, func(w http.ResponseWriter, r *http.Request) {})

	r.Route(static.DOMAIN_PATH, func(r chi.Router) {

		r.Route(static.USER_ID_PATH, func(r chi.Router) {
			r.Use(s.contextUid)
			r.Post(static.CREATE_PATH, s.creator.Create())

			r.Route(static.QUEUE_ID_PATH, func(r chi.Router) {
				r.Use(s.contextQid)
				r.Get(static.STATUS_PATH, s.checker.Check())
			})
		})
	})

	return r
}

func (s Server) contextUid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, string(static.USER_ID))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), static.USER_ID, uid)))
	})
}

func (s Server) contextQid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qid := chi.URLParam(r, string(static.QUEUE_ID))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), static.QUEUE_ID, qid)))
	})
}
