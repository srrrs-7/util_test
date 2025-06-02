package service

import (
	"concurrency/controller/request"
	"concurrency/domain"
	"context"
	"fmt"
	"time"
)

type ConcurrentService struct {
	queueRepo QueueRepo
	cacheRepo CacheRepo
}

func (s *ConcurrentService) Create(ctx context.Context, req request.CreateReq) (*domain.UserStatus, error) {
	user, err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	qid, err := s.queueRepo.Enqueue(ctx, user.ID.String())
	if err != nil {
		return nil, err
	}

	cacheKey := s.cacheRepo.MakeKey(ctx, "concurrent:", qid.String())
	if err := s.cacheRepo.Set(
		ctx,
		cacheKey,
		domain.StatusPending.String(),
		5*time.Minute,
	); err != nil {
		return nil, err
	}

	return &domain.UserStatus{
		QueueID: qid,
		UserID:  user.ID,
		Status:  domain.StatusPending,
	}, nil
}

func (s *ConcurrentService) Check(ctx context.Context, req request.CheckReq) (*domain.UserStatus, error) {
	qid, err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	cacheKey := s.cacheRepo.MakeKey(ctx, "concurrent:", qid.String())
	status, err := s.cacheRepo.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	switch status {
	case domain.StatusPending.String():
		return &domain.UserStatus{
			QueueID: qid,
			Status:  domain.StatusPending,
		}, nil
	case domain.StatusRunning.String():
		return &domain.UserStatus{
			QueueID: qid,
			Status:  domain.StatusRunning,
		}, nil
	case domain.StatusCompleted.String():
		return &domain.UserStatus{
			QueueID: qid,
			Status:  domain.StatusCompleted,
		}, nil
	case domain.StatusFailed.String():
		return &domain.UserStatus{
			QueueID: qid,
			Status:  domain.StatusFailed,
		}, nil
	default:
		return nil, fmt.Errorf("unknown status: %s", status)
	}
}
