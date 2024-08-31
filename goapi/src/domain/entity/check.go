package entity

import (
	"api/util/static"
	"fmt"
)

type Status string
type UserId string

type CheckStatusEnt struct {
	Id     QueueId `json:"id"`
	UserId UserId  `json:"userId"`
	Status Status  `json:"status"`
}

func (s Status) State() Status {
	switch s {
	case static.PENDING:
		return static.PENDING
	case static.DONE:
		return static.DONE
	case static.ERROR:
		return static.ERROR
	case static.RUNNING:
		return static.RUNNING
	default:
		return Status(fmt.Sprintf("invalid status: %s", s))
	}
}
