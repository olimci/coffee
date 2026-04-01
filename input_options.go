package coffee

type inputOptions struct {
	validate    func(string) error
	charLimit   int
	placeholder string
	suggestions []string
	width       int
	value       string
	valueSet    bool
	discard     bool
}

type InputOption func(*inputOptions)

func defaultInputOptions() *inputOptions {
	return &inputOptions{}
}

func (o *inputOptions) apply(opts ...InputOption) *inputOptions {
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func WithInputValidate(validate func(string) error) InputOption {
	return func(o *inputOptions) {
		o.validate = validate
	}
}

func WithInputCharLimit(limit int) InputOption {
	return func(o *inputOptions) {
		o.charLimit = limit
	}
}

func WithInputPlaceholder(placeholder string) InputOption {
	return func(o *inputOptions) {
		o.placeholder = placeholder
	}
}

func WithInputSuggestions(suggestions []string) InputOption {
	return func(o *inputOptions) {
		o.suggestions = suggestions
	}
}

func WithInputWidth(width int) InputOption {
	return func(o *inputOptions) {
		o.width = width
	}
}

func WithInputValue(value string) InputOption {
	return func(o *inputOptions) {
		o.value = value
		o.valueSet = true
	}
}

func WithInputDiscardSubmitted() InputOption {
	return func(o *inputOptions) {
		o.discard = true
	}
}

func configuredInput(opts *inputOptions) *Input {
	input := NewInput().
		WithValidate(opts.validate).
		WithCharLimit(opts.charLimit).
		WithPlaceholder(opts.placeholder).
		WithSuggestions(opts.suggestions).
		WithWidth(opts.width).
		WithDiscardSubmitted(opts.discard)

	if opts.valueSet {
		input.WithValue(opts.value)
	}

	return input
}
