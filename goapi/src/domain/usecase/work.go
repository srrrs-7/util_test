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
	qCh := make(chan *model.QueueModel[model.QueueModel[request.Params]])
	errCh := make(chan error)
	timer := time.NewTicker(1 * time.Second)

	defer func() {
		close(qCh)
		close(errCh)
		timer.Stop()
	}()

	go u.receiveQueue(ctx, qCh, errCh)

	for {
		if u.thread.cnt >= static.MAX_THREAD_CNT {
			slog.Info("max thread cnt reached", "count", u.thread.cnt)
			time.Sleep(1 * time.Second)
			continue
		}

		select {
		case <-timer.C:
			slog.Info("ticker", "thread count", u.thread.cnt)
		case <-ctx.Done():
			slog.Info("context done")
		case err := <-errCh:
			slog.Error("work error", "error", err.Error())
		case p := <-qCh:
			slog.Info("concurrency work", "param", p)
			go u.concurrencyWork(ctx, p)
		}
	}
}

func (u *WorkerUseCase) concurrencyWork(ctx context.Context, p *model.QueueModel[model.QueueModel[request.Params]]) {
	u.thread.increment()
	defer u.thread.decrement()

	status := entity.CheckStatusEnt{
		Id:     p.Entity().Id,
		UserId: entity.UserId(1),
		Status: entity.Status(static.PENDING),
	}

	if err := u.runWork(p); err != nil {
		status.Status = static.ERROR
		if err := u.cache.Set(ctx, p.Entity().Id, status); err != nil {
			slog.Error("concurrency work set status error", "error", err.Error())
			return
		}
		slog.Error("concurrency work error", "error", err.Error())
		return
	}

	status.Status = static.DONE
	if err := u.cache.Set(ctx, p.Entity().Id, status); err != nil {
		slog.Error("concurrency work set status error", "error", err.Error())
		return
	}

	ctx.Done()
}

func (u WorkerUseCase) receiveQueue(ctx context.Context, qCh chan<- *model.QueueModel[model.QueueModel[request.Params]], errCh chan<- error) {
	for {
		p, err := u.queue.DeQueue(ctx)
		if err != nil {
			errCh <- err
		}

		qCh <- p
	}
}

func (u WorkerUseCase) runWork(p *model.QueueModel[model.QueueModel[request.Params]]) error {
	fmt.Printf("work param: %v", p)
	time.Sleep(3 * time.Second)
	return nil
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
