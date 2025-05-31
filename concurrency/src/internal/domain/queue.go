package domain

import "fmt"

type QueueID string

func (q QueueID) String() string {
	if q == "" {
		return "unknown"
	}
	return string(q)
}

func (q QueueID) validate() error {
	if q == "" {
		return fmt.Errorf("queue ID cannot be empty")
	}
	return nil
}
