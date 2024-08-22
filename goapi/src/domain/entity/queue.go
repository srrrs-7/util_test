package entity

type QueueId string
type ReceiptId string

type QueueEnt[T any] struct {
	Id        QueueId
	Body      T
	ReceiptId ReceiptId
}

func (qid QueueId) String() string {
	return string(qid)
}

func (rid ReceiptId) String() string {
	return string(rid)
}
