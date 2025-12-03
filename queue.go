package main

import (
	"errors"
	"sync"
)

type Queue[T any] struct {
	queue []T
	mu    sync.Mutex
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}

func (q *Queue[T]) Enqueue(items ...T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, items...)
}

func (q *Queue[T]) Dequeue() (T, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		var none T
		return none, errors.New("queue is empty")
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item, nil
}

func (q *Queue[T]) Empty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue) == 0
}
