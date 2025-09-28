package util

import (
	"sync"
)

type WorkQueue[T any] struct {
	mu     sync.Mutex
	queue  chan T
	closed bool
}

func NewWorkQueue[T any](buffer int) *WorkQueue[T] {
	return &WorkQueue[T]{
		queue: make(chan T, buffer),
	}
}

func (wq *WorkQueue[T]) Add(item T) bool {
	wq.mu.Lock()
	defer wq.mu.Unlock()

	if wq.closed {
		return false
	}
	wq.queue <- item
	return true
}

func (wq *WorkQueue[T]) Get() (T, bool) {
	item, ok := <-wq.queue
	var zero T
	if !ok {
		return zero, false
	}
	return item, true
}

func (wq *WorkQueue[T]) Close() {
	wq.mu.Lock()
	defer wq.mu.Unlock()

	if !wq.closed {
		wq.closed = true
		close(wq.queue)
	}
}
