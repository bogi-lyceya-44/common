package worker

import (
	"context"
	"sync"
	"time"
)

type ProcessFunc[T any] func(context.Context, []T)

type WorkerPool[T any] struct {
	input chan T

	workerCount  int
	batchSize    int
	flushTimeout time.Duration
	process      ProcessFunc[T]

	cancel context.CancelFunc

	once *sync.Once
	wg   *sync.WaitGroup
}

func New[T any](
	inputBufferSize int,
	workerCount int,
	batchSize int,
	flushTimeout time.Duration,
	process ProcessFunc[T],
) *WorkerPool[T] {
	return &WorkerPool[T]{
		input:        make(chan T, inputBufferSize),
		workerCount:  workerCount,
		flushTimeout: flushTimeout,
		batchSize:    batchSize,
		process:      process,
		once:         &sync.Once{},
		wg:           &sync.WaitGroup{},
	}
}

func (w *WorkerPool[T]) Start(ctx context.Context) {
	go w.start(ctx)
}

func (w *WorkerPool[T]) start(ctx context.Context) {
	w.once.Do(func() {
		ctx, cancel := context.WithCancel(ctx)
		w.cancel = cancel

		for range w.workerCount {
			w.wg.Add(1)

			go func() {
				defer w.wg.Done()
				w.work(ctx)
			}()
		}
	})
}

func (w *WorkerPool[T]) work(ctx context.Context) {
	ticker := time.NewTicker(w.flushTimeout)
	defer ticker.Stop()

	batch := make([]T, 0, w.batchSize)

	flush := func(flushCtx context.Context) {
		if len(batch) == 0 {
			return
		}

		w.process(flushCtx, batch)
		batch = batch[:0]
	}

	for {
		select {
		case <-ticker.C:
			flush(ctx)
		case item, ok := <-w.input:
			if !ok {
				flush(context.Background())
				return
			}

			ticker.Reset(w.flushTimeout)
			batch = append(batch, item)

			if len(batch) >= w.batchSize {
				flush(ctx)
			}
		case <-ctx.Done():
			// for graceful shutdown:
			// flush the leftovers with context.Background()
			// since the original context is done
			flush(context.Background())
			return
		}
	}
}

func (w *WorkerPool[T]) Stop() {
	if w.cancel != nil {
		w.cancel()
	}

	w.wg.Wait()

	// create a new instance of once
	// so we'll be able to start worker again
	w.once = &sync.Once{}
}

func (w *WorkerPool[T]) Send(item T) {
	w.input <- item
}
