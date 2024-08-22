package repository

import (
	"api/domain/entity"
	"api/driver/model"
	"context"
	"encoding/json"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueRepo[T any] struct {
	client *sqs.Client
	url    string
}

func NewQueueRepo[T any](client *sqs.Client, url string) QueueRepo[T] {
	return QueueRepo[T]{client, url}
}

func (q QueueRepo[T]) EnQueue(ctx context.Context, msg string) (entity.QueueId, error) {
	res, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    (*string)(unsafe.Pointer(&q.url)),
		MessageBody: (*string)(unsafe.Pointer(&msg)),
	})
	if err != nil {
		return "", err
	}

	return entity.QueueId(*(*string)(unsafe.Pointer(res.MessageId))), nil
}

func (q QueueRepo[T]) DeQueue(ctx context.Context) (*model.QueueModel[T], error) {
	res, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            (*string)(unsafe.Pointer(&q.url)),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     1,
		VisibilityTimeout:   3,
	})
	if err != nil {
		return nil, err
	}

	var j T
	if err = json.Unmarshal([]byte(*res.Messages[0].Body), &j); err != nil {
		return nil, err
	}

	return &model.QueueModel[T]{
		Id:        *res.Messages[0].MessageId,
		Body:      j,
		ReceiptId: *res.Messages[0].ReceiptHandle,
	}, nil
}

func (q QueueRepo[T]) DelQueue(ctx context.Context, id string) error {
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      (*string)(unsafe.Pointer(&q.url)),
		ReceiptHandle: (*string)(unsafe.Pointer(&id)),
	})
	if err != nil {
		return err
	}

	return nil
}
