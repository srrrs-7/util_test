package handle

import (
	"api/domain/usecase"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type UserId string

const (
	HEALTH_PATH  = "/health"
	DOMAIN_PATH  = "/domain/v1"
	USER_ID_PATH = "/user/{userId}"
	CREATE_PATH  = "/create"
	STATUS_PATH  = "/status"

	USER_ID UserId = "userId"
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

	r.Get(HEALTH_PATH, func(w http.ResponseWriter, r *http.Request) {})

	r.Route(DOMAIN_PATH, func(r chi.Router) {
		r.Route(USER_ID_PATH, func(r chi.Router) {
			r.Use(s.contextUid)
			r.Group(func(r chi.Router) {
				r.Post(CREATE_PATH, s.creator.Create())
				r.Get(STATUS_PATH, s.checker.Check())
			})
		})

	})

	return r
}

func (s Server) contextUid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := chi.URLParam(r, string(USER_ID))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), USER_ID, uid)))
	})
}
