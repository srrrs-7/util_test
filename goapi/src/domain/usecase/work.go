package usecase

import (
	"api/domain"
	"api/domain/entity"
	"api/driver/model"
	"api/handle/request"
	"api/util/static"
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

type thread struct {
	mu  *sync.Mutex
	cnt int
}

type WorkerUseCase struct {
	queue  domain.Queuer[model.QueueModel[request.Params]]
	cache  domain.Cacher[entity.CheckStatusEnt]
	thread *thread
}

func NewWorkerUseCase(
	m *sync.Mutex,
	q domain.Queuer[model.QueueModel[request.Params]],
	c domain.Cacher[entity.CheckStatusEnt],
) *WorkerUseCase {
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

	defer func() {
		close(doneCh)
		close(errCh)
	}()

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
		if u.thread.cnt >= static.MAX_THREAD_CNT {
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
