package queue

import (
	"claude/driver/grpc/model"
	"fmt"
	"sync/atomic"
)

var queue atomic.Pointer[[]model.Prompt]

func NewQueue() {
	prompts := make([]model.Prompt, 0)
	queue.Store(&prompts)
}

func Enqueue(prompt model.Prompt) error {
	prompts := queue.Load()
	if prompts == nil {
		return fmt.Errorf("queue is not initialized")
	}

	*prompts = append(*prompts, prompt)

	return nil
}

func Dequeue() (model.Prompt, error) {
	prompts := queue.Load()
	if prompts == nil {
		return "", fmt.Errorf("queue is not initialized")
	}

	if len(*prompts) == 0 {
		return "", fmt.Errorf("no prompts in the queue")
	}

	prompt := (*prompts)[0]
	*prompts = (*prompts)[1:]

	return prompt, nil
}

func QueueStatus() (int, error) {
	prompts := queue.Load()
	if prompts == nil {
		return 0, fmt.Errorf("queue is not initialized")
	}

	l := len(*prompts)
	if l == 0 {
		return 0, fmt.Errorf("queue is empty")
	}

	return l, nil
}
