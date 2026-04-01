package coffee

import (
	"context"
	"errors"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

var ErrNonInteractive = errors.New("coffee requires an interactive terminal")

type Coffee struct {
	ctx    context.Context
	cancel context.CancelCauseFunc

	program *tea.Program
}

func (c *Coffee) Context() context.Context {
	return c.ctx
}

func Do(f func(ctx context.Context, c *Coffee) error, opts ...Option) error {
	o := defaultOptions().apply(opts...)
	c := new(o)
	defer func() {
		c.program.Quit()
		c.program.Wait()
	}()

	go func() {
		if _, err := c.program.Run(); err != nil {
			if isNonInteractive(err) {
				c.cancel(fmt.Errorf("%w: %w", ErrNonInteractive, err))
				return
			}
			c.cancel(err)
		} else {
			c.cancel(nil)
		}
	}()

	return cleanCancel(c.ctx, f(c.ctx, c))
}

func new(o *options) *Coffee {
	ctx, cancel := context.WithCancelCause(o.ctx)

	opts := []tea.ProgramOption{tea.WithContext(ctx)}
	if o.altScreen {
		opts = append(opts, tea.WithAltScreen())
	}

	return &Coffee{
		ctx:     ctx,
		cancel:  cancel,
		program: tea.NewProgram(newModel(o), opts...),
	}
}
