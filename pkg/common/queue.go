package common

import (
	"errors"
)

// @TODO: Make it auto resizable, @see https://github.com/gammazero/deque
type Queue[T any] struct {
	capacity int
	q        chan T
}

type FifoQueue interface {
	Append()
	Pop()
}

func NewQueue[T any](capacity int) *Queue[T] {
	return &Queue[T]{
		capacity: capacity,
		q:        make(chan T, capacity),
	}
}

func (q *Queue[T]) Append(item T) error {
	if len(q.q) < int(q.capacity) {
		q.q <- item
		return nil
	}

	return errors.New("queue is full")
}

func (q *Queue[T]) Pop() (*T, error) {
	if len(q.q) > 0 {
		item := <-q.q
		return &item, nil
	}

	return nil, errors.New("queue is empty")
}

func (q *Queue[T]) Empty() bool {
	return len(q.q) == 0
}
