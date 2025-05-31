package router

import (
	"concurrency/internal/controller/request"
	"concurrency/internal/controller/response"
	"concurrency/internal/service"
	"concurrency/pkg/utilhttp"
	"log/slog"
	"net/http"
)

type Concurrent struct {
	service *service.ConcurrentService
}

func (c *Concurrent) Create(w http.ResponseWriter, r *http.Request) {
	req, err := utilhttp.GetRequestJsonBody[request.CreateReq](r)
	if err != nil {
		slog.Error("Failed to parse request body", "error", err)
		utilhttp.ResponseBadRequest(w, err.Error())
		return
	}

	status, err := c.service.Create(r.Context(), req)
	if err != nil {
		slog.Error("Failed to create concurrent request", "error", err)
		utilhttp.ResponseInternalServerError(w, err.Error())
		return
	}

	utilhttp.ResponseAccepted(w, response.ConcurrentResp{
		ID:     status.QueueID.String(),
		Status: status.Status.String(),
	})
}

func (c *Concurrent) Check(w http.ResponseWriter, r *http.Request) {
	req, err := utilhttp.GetRequestJsonBody[request.CheckReq](r)
	if err != nil {
		slog.Error("Failed to parse request body", "error", err)
		utilhttp.ResponseBadRequest(w, err.Error())
		return
	}

	status, err := c.service.Check(r.Context(), req)
	if err != nil {
		slog.Error("Failed to check concurrent request", "error", err)
		utilhttp.ResponseInternalServerError(w, err.Error())
		return
	}

	utilhttp.ResponseOk(w, response.ConcurrentResp{
		ID:     status.QueueID.String(),
		Status: status.Status.String(),
	})
}
