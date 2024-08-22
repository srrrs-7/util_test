package model

import "api/domain/entity"

type QueueModel[T any] struct {
	Id        string
	Body      T
	ReceiptId string
}

func (q QueueModel[T]) Entity() entity.QueueEnt[T] {
	return entity.QueueEnt[T]{
		Id:        entity.QueueId(q.Id),
		Body:      q.Body,
		ReceiptId: entity.ReceiptId(q.ReceiptId),
	}
}
