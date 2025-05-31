package queue

import (
	"concurrency/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Worker interface {
	worker(context.Context, string, string) error
}

type SQSHandler struct {
	client   *sqs.Client
	queueURL string
	timeout  time.Duration
	work     Worker
}

func NewSQSHandler(queueUrl string, timeout time.Duration, w Worker) (*SQSHandler, error) {
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

	return &SQSHandler{
		client:   client,
		queueURL: *result.QueueUrl,
		timeout:  timeout,
		work:     w,
	}, nil
}

func (h *SQSHandler) Enqueue(ctx context.Context, messageBody string) (domain.QueueID, error) {
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

	return domain.QueueID(*result.MessageId), nil
}

func (h *SQSHandler) Dequeue(ctx context.Context) error {
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
		time.Sleep(3 * time.Second)
		return fmt.Errorf("no messages in the queue")
	}

	msg := result.Messages[0]

	ticker := time.NewTicker(1 * time.Second)
	go func(ticker *time.Ticker) {
		for range ticker.C {
			if err := h.extendVisibility(ctx, *msg.ReceiptHandle, int32(h.timeout.Seconds())); err != nil {
				slog.Error("extending message visibility", "error", err.Error())
				return
			}
		}
	}(ticker)

	defer func() {
		ticker.Stop()
		if err := h.deleteMessage(ctx, *result.Messages[0].ReceiptHandle); err != nil {
			slog.Error("deleting message", "error", err.Error())
			return
		}
	}()

	if err := h.work.worker(ctx, *msg.MessageId, *msg.Body); err != nil {
		return fmt.Errorf("processing message: %w", err)
	}

	return nil
}

func (h *SQSHandler) deleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := h.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(h.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	if err != nil {
		return fmt.Errorf("deleting message: %w", err)
	}

	return nil
}

func (h *SQSHandler) extendVisibility(ctx context.Context, receiptHandle string, visibilityTimeout int32) error {
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
