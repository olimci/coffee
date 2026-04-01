package promise

import (
	"context"
)

func New[T any]() (Promise[T], Resolver[T]) {
	ch := make(chan result[T])
	return Promise[T]{ch}, Resolver[T]{ch}
}

type Promise[T any] struct {
	ch <-chan result[T]
}

func (p Promise[T]) Await(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		return *new(T), ctx.Err()
	case res := <-p.ch:
		return res.value, res.err
	}
}

type Resolver[T any] struct {
	ch chan<- result[T]
}

func (r Resolver[T]) resolve(res result[T]) {
	r.ch <- res
}

func (r Resolver[T]) Ok(value T) {
	r.resolve(result[T]{value: value})
}

func (r Resolver[T]) Err(err error) {
	r.resolve(result[T]{err: err})
}

type result[T any] struct {
	value T
	err   error
}
