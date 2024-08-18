package usecase

import (
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
	thread *thread
}

func NewWorkerUseCase(m *sync.Mutex) *WorkerUseCase {
	return &WorkerUseCase{
		&thread{
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
			time.Sleep(100 * time.Microsecond)
			continue
		}
		// dequeue

		// set status

		u.threadCntUp()
		// work logic
		u.threadCntDown()

		// set new status

		doneCh <- struct{}{}
		errCh <- nil
		ctx.Done()
	}
}

func (u *WorkerUseCase) threadCntUp() {
	u.thread.mu.Lock()
	defer u.thread.mu.Unlock()
	u.thread.cnt++
}

func (u *WorkerUseCase) threadCntDown() {
	u.thread.mu.Lock()
	defer u.thread.mu.Unlock()
	u.thread.cnt--
}
