package coffee

type logOptions struct {
	wrap   bool
	indent string
}

// LogOption configures how a log message is rendered.
type LogOption func(*logOptions)

func defaultLogOptions() logOptions {
	return logOptions{}
}

// WithWrap reflows a log message to the current viewport width on each render.
func WithWrap() LogOption {
	return func(o *logOptions) {
		o.wrap = true
	}
}

// WithIndent prefixes each rendered log line with the provided indent.
func WithIndent(indent string) LogOption {
	return func(o *logOptions) {
		o.indent = indent
	}
}

func applyLogOptions(opts ...LogOption) logOptions {
	out := defaultLogOptions()
	for _, opt := range opts {
		opt(&out)
	}
	return out
}
