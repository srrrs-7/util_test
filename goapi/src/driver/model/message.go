package model

type QueueId string

type QueueBody struct {
	Msg string `json:"message"`
}
