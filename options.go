package coffee

import "context"

func defaultOptions() *options {
	return &options{
		ctx:       context.Background(),
		altScreen: false,
	}
}

type options struct {
	ctx       context.Context
	altScreen bool
}

func (o *options) apply(opts ...Option) *options {
	for _, opt := range opts {
		opt(o)
	}

	return o
}

type Option func(*options)

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithAltScreen() Option {
	return func(o *options) {
		o.altScreen = true
	}
}
