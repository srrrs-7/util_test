package entity

import "api/handle/response"

type QueueId string
type Status string

type CheckStatusEnt struct {
	Id     QueueId
	Status Status
}

func (e *CheckStatusEnt) Response() response.StatusRes {
	return response.StatusRes{
		Id:     string(e.Id),
		Status: string(e.Status),
	}
}
