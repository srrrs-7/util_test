package usecase

import (
	"api/util/utilhttp"
	"net/http"
)

type CreateUseCase struct{}

// API request -> send queue -> set redis -> API response
func (u CreateUseCase) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilhttp.ResponseOk(w, nil)
	}
}
