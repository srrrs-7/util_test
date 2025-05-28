package service

import "concurrency/internal/domain"

type Work struct{}

func (c *Check) Work(id string) (*domain.UserStatus, error) {
	return nil, nil
}
