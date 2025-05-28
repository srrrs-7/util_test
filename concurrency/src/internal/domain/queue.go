package domain

import "concurrency/pkg/utilerror"

type QueueID string

func (q QueueID) String() string {
	if q == "" {
		return "unknown"
	}
	return string(q)
}

func (q QueueID) validate() error {
	if q == "" {
		return utilerror.ErrInvalidQueueID
	}
	return nil
}
