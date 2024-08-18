package usecase

import (
	"api/domain"
	"api/driver/model"
	"api/handle/request"
	"api/handle/response"
	"api/util/utilhttp"
	"log/slog"
	"net/http"
)

type CreateUseCase struct {
	queue domain.Queuer[model.QueueModel]
	cache domain.Cacher[model.CacheModel]
}

func NewCreateUseCase(queue domain.Queuer[model.QueueModel], cache domain.Cacher[model.CacheModel]) CreateUseCase {
	return CreateUseCase{queue, cache}
}

// API request -> send queue -> set redis -> API response
func (u CreateUseCase) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j, err := utilhttp.Json(utilhttp.RequestBody[request.Params](r))
		if err != nil {
			slog.Error("request json error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		qId, err := u.queue.EnQueue(r.Context(), string(j))
		if err != nil {
			slog.Error("enqueue error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		userId := utilhttp.RequestUrlParam(r, "userId")
		if err = u.cache.Set(r.Context(), userId, qId); err != nil {
			slog.Error("set cache error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		j, err = utilhttp.Json(response.StatusRes{Id: qId, Status: response.PENDING})
		if err != nil {
			slog.Error("response json error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		utilhttp.ResponseOk(w, j)
	}
}
