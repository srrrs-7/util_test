package grpc

import (
	"claude/driver/grpc/model"
	"claude/util/queue"
	"context"
	"fmt"
)

type Enqueue struct {
	UnimplementedEnqueueServiceServer
}

func (e *Enqueue) Enqueue(ctx context.Context, req *EnqueueRequest) (*EnqueueResponse, error) {
	if err := queue.Enqueue(model.Prompt(req.Prompt)); err != nil {
		return nil, fmt.Errorf("failed to enqueue prompt: %w", err)
	}

	return &EnqueueResponse{
		Message: "Message enqueued successfully",
	}, nil
}

type Dequeue struct {
	UnimplementedDequeueServiceServer
}

func (d *Dequeue) Dequeue(ctx context.Context, req *DequeueRequest) (*DequeueResponse, error) {
	prompt, err := queue.Dequeue()
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue prompt: %w", err)
	}

	return &DequeueResponse{
		Prompt: string(prompt),
	}, nil
}

type QueueStatus struct {
	UnimplementedQueueStatusServiceServer
}

func (qs *QueueStatus) GetQueueStatus(ctx context.Context, req *QueueStatusRequest) (*QueueStatusResponse, error) {
	l, err := queue.QueueStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get queue status: %w", err)
	}
	if l == 0 {
		return &QueueStatusResponse{
			Status: "Queue is empty",
		}, nil
	}

	return &QueueStatusResponse{
		Status: fmt.Sprintf("Queue has %d prompts", l),
	}, nil
}
