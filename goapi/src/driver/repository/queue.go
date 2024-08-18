package repository

import (
	"context"
	"encoding/json"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueRepo[T any] struct {
	client *sqs.Client
	url    string
}

func NewQueueRepo[T any](c *sqs.Client, url string) QueueRepo[T] {
	return QueueRepo[T]{
		client: c,
		url:    url,
	}
}

func (q QueueRepo[T]) Send(ctx context.Context, msg string) (string, error) {
	res, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    (*string)(unsafe.Pointer(&q.url)),
		MessageBody: (*string)(unsafe.Pointer(&msg)),
	})
	if err != nil {
		return "", err
	}

	return *(*string)(unsafe.Pointer(res.MessageId)), nil

}

func (q QueueRepo[T]) Receive(ctx context.Context) (*T, error) {
	res, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            (*string)(unsafe.Pointer(&q.url)),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     1,
		VisibilityTimeout:   1,
	})
	if err != nil {
		return nil, err
	}

	var j T
	if err = json.Unmarshal([]byte(*res.Messages[0].Body), &j); err != nil {
		return nil, err
	}

	return &j, nil
}

func (q QueueRepo[T]) Delete(ctx context.Context, id string) error {
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      (*string)(unsafe.Pointer(&q.url)),
		ReceiptHandle: (*string)(unsafe.Pointer(&id)),
	})
	if err != nil {
		return err
	}

	return nil
}
