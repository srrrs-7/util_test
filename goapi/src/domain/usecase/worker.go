package usecase

import (
	"api/domain"
	"api/driver/model"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

const (
	MAX_THREAD_CNT = 100
)

type thread struct {
	mu  *sync.Mutex
	cnt int
}

type WorkerUseCase struct {
	queue  domain.Queuer[model.QueueModel]
	cache  domain.Cacher[model.CacheModel]
	thread *thread
}

func NewWorkerUseCase(m *sync.Mutex, q domain.Queuer[model.QueueModel], c domain.Cacher[model.CacheModel]) *WorkerUseCase {
	return &WorkerUseCase{
		queue: q,
		cache: c,
		thread: &thread{
			mu:  m,
			cnt: 0,
		},
	}
}

// receive queue -> set new status -> processing -> delete queue -> set new status
func (u *WorkerUseCase) Work(ctx context.Context) {
	doneCh := make(chan struct{})
	errCh := make(chan error)

	defer close(doneCh)
	defer close(errCh)

	go u.concurrencyWork(ctx, doneCh, errCh)

	for {
		select {
		case s := <-doneCh:
			fmt.Println(s)

		case err := <-errCh:
			slog.Error("work concurrency error", "error", err.Error())
		}
	}
}

func (u *WorkerUseCase) concurrencyWork(
	ctx context.Context,
	doneCh chan<- struct{},
	errCh chan<- error,
) {
	for {
		if u.thread.cnt >= MAX_THREAD_CNT {
			slog.Info("max thread cnt reached", "count", u.thread.cnt)
			time.Sleep(100 * time.Microsecond)
			continue
		}
		// dequeue

		// set status

		u.thread.increment()
		// work logic
		u.thread.decrement()

		// set new status

		doneCh <- struct{}{}
		errCh <- nil
		ctx.Done()
	}
}

func (t *thread) increment() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cnt++
}

func (t *thread) decrement() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cnt--
}
