package usecase

import (
	"api/domain"
	"api/domain/entity"
	"api/driver/model"
	"api/handle/request"
	"api/handle/response"
	"api/util/static"
	"api/util/utilhttp"
	"log/slog"
	"net/http"
)

type CreateUseCase struct {
	queue domain.Queuer[model.QueueModel[request.Params]]
	cache domain.Cacher[entity.CheckStatusEnt]
}

func NewCreateUseCase(
	queue domain.Queuer[model.QueueModel[request.Params]],
	cache domain.Cacher[entity.CheckStatusEnt],
) CreateUseCase {
	return CreateUseCase{queue, cache}
}

// API request -> send queue -> set redis -> API response
func (u CreateUseCase) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		j, err := utilhttp.Json(utilhttp.RequestBody[request.Params](r))
		if err != nil {
			slog.Error("request json error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseBadRequest(w, j)
			return
		}

		qId, err := u.queue.EnQueue(r.Context(), string(j))
		if err != nil {
			slog.Error("enqueue error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		userId, err := utilhttp.RequestUrlParam[string](r, "userId")
		if err != nil {
			slog.Error("request url param error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseBadRequest(w, j)
			return
		}

		status := entity.CheckStatusEnt{
			Id:     qId,
			UserId: entity.UserId(userId),
			Status: entity.PENDING,
		}
		if err = u.cache.Set(r.Context(), qId, status); err != nil {
			slog.Error("set cache error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		j, err = utilhttp.Json(response.StatusRes{Id: qId.String(), Status: static.PENDING})
		if err != nil {
			slog.Error("response json error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		utilhttp.ResponseOk(w, j)
	}
}
