package repository

import (
	"context"
	"unsafe"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type QueueRepo struct {
	client *sqs.Client
	url    string
}

func NewQueueRepo(c *sqs.Client, url string) QueueRepo {
	return QueueRepo{
		client: c,
		url:    url,
	}
}

func (q QueueRepo) Send(ctx context.Context, msg string) (string, error) {
	res, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    (*string)(unsafe.Pointer(&q.url)),
		MessageBody: (*string)(unsafe.Pointer(&msg)),
	})
	if err != nil {
		return "", err
	}

	return *(*string)(unsafe.Pointer(res.MessageId)), nil

}

func (q QueueRepo) Receive(ctx context.Context) (string, error) {
	res, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            (*string)(unsafe.Pointer(&q.url)),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     1,
		VisibilityTimeout:   1,
	})
	if err != nil {
		return "", err
	}

	if len(res.Messages) == 0 {
		return "", nil
	}

	return *(*string)(unsafe.Pointer(res.Messages[0].Body)), nil
}

func (q QueueRepo) Delete(ctx context.Context, id string) error {
	_, err := q.client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
		QueueUrl:      (*string)(unsafe.Pointer(&q.url)),
		ReceiptHandle: (*string)(unsafe.Pointer(&id)),
	})
	if err != nil {
		return err
	}

	return nil
}
