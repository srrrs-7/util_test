package request

import (
	"concurrency/domain"
	"fmt"
)

type CreateReq struct {
	UserID string `json:"user_id"`
}

func (c *CreateReq) Validate() (*domain.User, error) {
	if c.UserID == "" {
		return nil, fmt.Errorf("user ID cannot be empty")
	}
	return &domain.User{ID: domain.UserID(c.UserID)}, nil
}

type CheckReq struct {
	QueueID string `json:"queue_id"`
}

func (c *CheckReq) Validate() (domain.QueueID, error) {
	if c.QueueID == "" {
		return "", fmt.Errorf("ID cannot be empty")
	}
	return domain.QueueID(c.QueueID), nil
}
