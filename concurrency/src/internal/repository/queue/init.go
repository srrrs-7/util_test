package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type WorkerFunc[T any] func(ctx context.Context, msgId string, msg T) error

type SQSHandler[T any] struct {
	client   *sqs.Client
	queueURL string
	timeout  time.Duration
	workFunc WorkerFunc[T]
}

func NewSQSHandler[T any](queueUrl string, timeout time.Duration, workFunc WorkerFunc[T]) (*SQSHandler[T], error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("loading AWS config: %w", err)
	}

	client := sqs.NewFromConfig(cfg)

	result, err := client.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueUrl),
	})
	if err != nil {
		return nil, fmt.Errorf("getting queue URL: %w", err)
	}

	return &SQSHandler[T]{
		client:   client,
		queueURL: *result.QueueUrl,
		timeout:  timeout,
		workFunc: workFunc,
	}, nil
}

func (h *SQSHandler[T]) Enqueue(ctx context.Context, messageBody T) (string, error) {
	messageBodyBytes, err := json.Marshal(messageBody)
	if err != nil {
		return "", fmt.Errorf("marshalling message: %w", err)
	}

	result, err := h.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(messageBodyBytes)),
		QueueUrl:    aws.String(h.queueURL),
	})
	if err != nil {
		return "", fmt.Errorf("sending message: %w", err)
	}

	return *result.MessageId, nil
}

func (h *SQSHandler[T]) Dequeue(ctx context.Context) error {
	result, err := h.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(h.queueURL),
		MaxNumberOfMessages: int32(1),
		WaitTimeSeconds:     int32(20),
	})
	if err != nil {
		return fmt.Errorf("receiving message: %w", err)
	}

	if len(result.Messages) == 0 {
		h.Dequeue(ctx)
		return fmt.Errorf("no messages in the queue")
	}

	msg := result.Messages[0]
	slog.Info("received message", "messageId", *msg.MessageId)

	ticker := time.NewTicker(4 * time.Minute)
	defer ticker.Stop()
	go func(ticker *time.Ticker) {
		for range ticker.C {
			if err := h.extendVisibility(ctx, *msg.ReceiptHandle, int32(h.timeout.Seconds())); err != nil {
				slog.Error("extending message visibility", "error", err.Error())
				return
			}
		}
	}(ticker)

	defer func() {
		if err := h.deleteMessage(ctx, *result.Messages[0].ReceiptHandle); err != nil {
			slog.Error("deleting message", "error", err.Error())
			return
		}
	}()

	var msgBody T
	if err := json.Unmarshal([]byte(*msg.Body), &msgBody); err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	if err := h.workFunc(ctx, *msg.MessageId, msgBody); err != nil {
		return fmt.Errorf("processing message: %w", err)
	}

	return nil
}

func (h *SQSHandler[T]) deleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := h.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(h.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}

	return nil
}

func (h *SQSHandler[T]) extendVisibility(ctx context.Context, receiptHandle string, visibilityTimeout int32) error {
	_, err := h.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(h.queueURL),
		ReceiptHandle:     aws.String(receiptHandle),
		VisibilityTimeout: int32(visibilityTimeout),
	})
	if err != nil {
		return fmt.Errorf("extending message visibility: %w", err)
	}

	return nil
}
