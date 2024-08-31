package repository

import (
	"api/domain/entity"
	"api/driver/model"
	"context"
	"encoding/json"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
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

func (q QueueRepo[T]) DeQueue(ctx context.Context) ([]*model.QueueModel[T], error) {
	res, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            (*string)(unsafe.Pointer(&q.url)),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     5,
		VisibilityTimeout:   5,
	})
	if err != nil {
		return nil, err
	}

	msgs := make([]*model.QueueModel[T], 0)
	entries := make([]types.DeleteMessageBatchRequestEntry, 0)
	for _, msg := range res.Messages {
		var j T
		if err = json.Unmarshal([]byte(*msg.Body), &j); err != nil {
			return nil, err
		}

		msgs = append(msgs, &model.QueueModel[T]{
			Id:        *msg.MessageId,
			Body:      j,
			ReceiptId: *msg.ReceiptHandle,
		})

		entries = append(entries, types.DeleteMessageBatchRequestEntry{
			Id:            msg.MessageId,
			ReceiptHandle: msg.ReceiptHandle,
		})
	}

	if len(entries) > 0 {
		if err = q.DelQueue(ctx, entries); err != nil {
			return nil, err
		}
	}

	return msgs, nil
}

func (q QueueRepo[T]) DelQueue(ctx context.Context, entries []types.DeleteMessageBatchRequestEntry) error {
	_, err := q.client.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		QueueUrl: (*string)(unsafe.Pointer(&q.url)),
		Entries:  entries,
	})
	if err != nil {
		return err
	}

	return nil
}
