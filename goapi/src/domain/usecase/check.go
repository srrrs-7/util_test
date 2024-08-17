package usecase

import (
	"api/util/utilhttp"
	"net/http"
)

type CheckUseCase struct{}

func (u *CheckUseCase) Check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utilhttp.ResponseOk(w, nil)
	}
}
