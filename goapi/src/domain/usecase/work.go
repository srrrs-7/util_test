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
)

type thread struct {
	cond *sync.Cond
	mu   *sync.Mutex
	cnt  int
}

type WorkerUseCase struct {
	queue  domain.Queuer[model.QueueModel[request.Params]]
	cache  domain.Cacher[entity.CheckStatusEnt]
	thread *thread
}

func NewWorkerUseCase(
	cond *sync.Cond,
	m *sync.Mutex,
	q domain.Queuer[model.QueueModel[request.Params]],
	c domain.Cacher[entity.CheckStatusEnt],
) *WorkerUseCase {
	return &WorkerUseCase{
		queue: q,
		cache: c,
		thread: &thread{
			cond: cond,
			mu:   m,
			cnt:  0,
		},
	}
}

// receive queue -> set new status -> processing -> delete queue -> set new status
func (u *WorkerUseCase) Work(ctx context.Context) {
	qCh := make(chan *model.QueueModel[model.QueueModel[request.Params]])
	errCh := make(chan error)

	defer func() {
		close(qCh)
		close(errCh)
	}()

	go u.receiveQueue(ctx, qCh, errCh)

	for {
		select {
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
		status.Status = entity.Status(static.ERROR)
		if err := u.cache.Set(ctx, p.Entity().Id, status); err != nil {
			slog.Error("concurrency work set status error", "error", err.Error())
			return
		}
		slog.Error("concurrency work error", "error", err.Error())
		return
	}

	status.Status = entity.Status(static.DONE)
	if err := u.cache.Set(ctx, p.Entity().Id, status); err != nil {
		slog.Error("concurrency work set status error", "error", err.Error())
		return
	}

	ctx.Done()
}

func (u WorkerUseCase) receiveQueue(ctx context.Context, qCh chan<- *model.QueueModel[model.QueueModel[request.Params]], errCh chan<- error) {
	for {
		msgs, err := u.queue.DeQueue(ctx)
		if err != nil {
			errCh <- err
		}

		for _, msg := range msgs {
			qCh <- msg
		}
	}
}

func (u WorkerUseCase) runWork(p *model.QueueModel[model.QueueModel[request.Params]]) error {
	fmt.Printf("work param: %v", p)
	return nil
}

func (t *thread) increment() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cnt++
	if t.cnt >= static.MAX_THREAD_CNT {
		t.cond.Wait()
	}
	t.cond.Broadcast()
}

func (t *thread) decrement() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.cnt--
}
