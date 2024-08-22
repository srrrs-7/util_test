package entity

const (
	PENDING Status = "PENDING"
	DONE    Status = "DONE"
	ERROR   Status = "ERROR"
)

type Status string
type UserId string

type CheckStatusEnt struct {
	Id     QueueId `json:"id"`
	UserId UserId  `json:"userId"`
	Status Status  `json:"status"`
}
