package service

import (
	"concurrency/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Work struct {
	cacheRepo CacheRepo
}

func (w *Work) Worker(ctx context.Context, queueId string, msg string) error {
	var user domain.User
	if err := json.Unmarshal([]byte(msg), &user); err != nil {
		return err
	}

	key := w.cacheRepo.MakeKey(ctx, "concurrent:", queueId)
	// status running
	if err := w.cacheRepo.Set(ctx, key, domain.StatusRunning.String(), 5*time.Minute); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	// process main logic
	if err := w.mainProcess(user); err != nil {
		// status failed
		if err := w.cacheRepo.Set(ctx, key, domain.StatusFailed.String(), 5*time.Minute); err != nil {
			return fmt.Errorf("failed to set cache: %w", err)
		}
		return fmt.Errorf("failed to process user: %w", err)
	}

	// status completed
	if err := w.cacheRepo.Set(ctx, key, domain.StatusCompleted.String(), 5*time.Minute); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

func (w *Work) mainProcess(user domain.User) error {
	fmt.Printf("Processing user: %s\n", user.ID.String())

	if err := user.Validate(); err != nil {
		return fmt.Errorf("invalid user: %w", err)
	}

	// Simulate some processing time
	time.Sleep(1 * time.Second)

	return nil
}
