package service

import "concurrency/internal/domain"

type Check struct{}

func (c *Check) Check(id string) (*domain.UserStatus, error) {
	return nil, nil
}
