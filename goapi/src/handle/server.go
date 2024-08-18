package handle

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	HEALTH_PATH  = "/health"
	DOMAIN_PATH  = "/domain/v1"
	USER_ID_PATH = "/user/{userId}"
	CREATE_PATH  = "/create"
	STATUS_PATH  = "/status"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (s Server) Routing() *chi.Mux {
	r := chi.NewMux()

	r.Get(HEALTH_PATH, func(w http.ResponseWriter, r *http.Request) {})

	r.Route(DOMAIN_PATH, func(r chi.Router) {
		r.Route(USER_ID_PATH, func(r chi.Router) {
			r.Use(r.Middlewares()...)
			r.Group(func(r chi.Router) {
				r.Post(CREATE_PATH, func(w http.ResponseWriter, r *http.Request) {})
				r.Get(STATUS_PATH, func(w http.ResponseWriter, r *http.Request) {})
			})
		})

	})

	return r
}
