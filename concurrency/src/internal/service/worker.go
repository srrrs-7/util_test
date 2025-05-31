package service

import "concurrency/internal/domain"

type Work struct{}

func (w *Work) Work(id string) (*domain.UserStatus, error) {
	return nil, nil
}
