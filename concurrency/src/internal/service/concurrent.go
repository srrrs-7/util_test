package service

import (
	"concurrency/internal/controller/request"
	"concurrency/internal/domain"
	"context"
	"fmt"
	"time"
)

type QueueRepo interface {
	Enqueue(ctx context.Context, messageBody string) (domain.QueueID, error)
	Dequeue(ctx context.Context) error
}

type CacheRepo interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	MakeKey(ctx context.Context, prefix, suffix string) string
}

type ConcurrentService struct {
	queueRepo QueueRepo
	cacheRepo CacheRepo
}

func (s *ConcurrentService) Create(ctx context.Context, req request.CreateReq) (*domain.UserStatus, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	id, err := s.queueRepo.Enqueue(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	cacheKey := s.cacheRepo.MakeKey(ctx, "concurrent:", id.String())
	if err := s.cacheRepo.Set(
		ctx,
		cacheKey,
		domain.StatusPending.String(),
		5*time.Minute,
	); err != nil {
		return nil, err
	}

	return &domain.UserStatus{
		QueueID: id,
		UserID:  domain.UserID(req.UserID),
		Status:  domain.StatusPending,
	}, nil
}

func (s *ConcurrentService) Check(ctx context.Context, req request.CheckReq) (*domain.UserStatus, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	cacheKey := s.cacheRepo.MakeKey(ctx, "concurrent:", req.ID)
	status, err := s.cacheRepo.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	switch status {
	case domain.StatusPending.String():
		return &domain.UserStatus{
			QueueID: domain.QueueID(req.ID),
			Status:  domain.StatusPending,
		}, nil
	case domain.StatusRunning.String():
		return &domain.UserStatus{
			QueueID: domain.QueueID(req.ID),
			Status:  domain.StatusRunning,
		}, nil
	case domain.StatusCompleted.String():
		return &domain.UserStatus{
			QueueID: domain.QueueID(req.ID),
			Status:  domain.StatusCompleted,
		}, nil
	case domain.StatusFailed.String():
		return &domain.UserStatus{
			QueueID: domain.QueueID(req.ID),
			Status:  domain.StatusFailed,
		}, nil
	default:
		return nil, fmt.Errorf("unknown status: %s", status)
	}
}
