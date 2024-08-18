package usecase

import (
	"api/domain"
	"api/driver/model"
	"api/handle/response"
	"api/util/utilhttp"
	"log/slog"
	"net/http"
)

type CheckUseCase struct {
	cache domain.Cacher[model.CacheModel]
}

func NewCheckUseCase(cache domain.Cacher[model.CacheModel]) CheckUseCase {
	return CheckUseCase{cache}
}

// API request -> get status -> API response
func (u CheckUseCase) Check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := utilhttp.RequestUrlParam(r, "userId")

		status, err := u.cache.Get(r.Context(), userId)
		if err != nil {
			slog.Error("get cache error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		j, err := utilhttp.Json(response.StatusRes{Id: userId, Status: string(status.Value)})
		if err != nil {
			slog.Error("response json error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		utilhttp.ResponseOk(w, j)
	}
}
