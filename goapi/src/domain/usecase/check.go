package usecase

import (
	"api/domain"
	"api/domain/entity"
	"api/handle/response"
	"api/util/static"
	"api/util/utilhttp"
	"log/slog"
	"net/http"
)

type CheckUseCase struct {
	cache domain.Cacher[entity.CheckStatusEnt]
}

func NewCheckUseCase(cache domain.Cacher[entity.CheckStatusEnt]) CheckUseCase {
	return CheckUseCase{cache}
}

// API request -> get status -> API response
func (u CheckUseCase) Check() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qid, err := utilhttp.RequestUrlParam[string](r, static.QUEUE_ID)
		if err != nil {
			slog.Error("request url param error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseBadRequest(w, j)
			return
		}

		state, err := u.cache.Get(r.Context(), entity.QueueId(qid))
		if err != nil {
			slog.Error("get cache error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		j, err := utilhttp.Json(response.StatusRes{Id: qid, Status: string(state.Status)})
		if err != nil {
			slog.Error("response json error", "error", err.Error())
			j, _ := utilhttp.Json(response.ErrorRes{Msg: err.Error()})
			utilhttp.ResponseInternalServerError(w, j)
			return
		}

		slog.Info("check status ok", "queueId", qid, "status", state.Status)
		utilhttp.ResponseOk(w, j)
	}
}
