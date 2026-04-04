package coffee

import (
	"fmt"
	"slices"

	"github.com/olimci/coffee/internal/promise"
)

type submodelOptions struct {
	section Section
	behind  bool
}

func defaultSubmodelOptions() submodelOptions {
	return submodelOptions{
		section: SectionBody,
	}
}

type SubmodelOption func(*submodelOptions)

func WithSection(section Section) SubmodelOption {
	return func(o *submodelOptions) {
		o.section = section
	}
}

func WithFocusBehind() SubmodelOption {
	return func(o *submodelOptions) {
		o.behind = true
	}
}

func (c *Coffee) Log(message string, opts ...LogOption) error {
	return c.send(msgLog{
		message: message,
		section: SectionBody,
		opts:    applyLogOptions(opts...),
	})
}

func (c *Coffee) Logf(format string, args ...any) error {
	return c.Log(fmt.Sprintf(format, args...))
}

func (c *Coffee) LogHeader(message string, opts ...LogOption) error {
	return c.send(msgLog{
		message: message,
		section: SectionHeader,
		opts:    applyLogOptions(opts...),
	})
}

func (c *Coffee) LogHeaderf(format string, args ...any) error {
	return c.LogHeader(fmt.Sprintf(format, args...))
}

func (c *Coffee) LogFooter(message string, opts ...LogOption) error {
	return c.send(msgLog{
		message: message,
		section: SectionFooter,
		opts:    applyLogOptions(opts...),
	})
}

func (c *Coffee) LogFooterf(format string, args ...any) error {
	return c.LogFooter(fmt.Sprintf(format, args...))
}

func (c *Coffee) Clear() error {
	return c.send(msgClear{section: SectionBody})
}

func (c *Coffee) ClearHeader() error {
	return c.send(msgClear{section: SectionHeader})
}

func (c *Coffee) ClearFooter() error {
	return c.send(msgClear{section: SectionFooter})
}

func (c *Coffee) SetWindowTitle(title string) error {
	return c.send(msgWindowTitle{title: title})
}

func (c *Coffee) AddSubmodel(submodel Submodel, opts ...SubmodelOption) error {
	o := defaultSubmodelOptions()
	for _, opt := range opts {
		opt(&o)
	}

	return c.send(msgSubmodel{
		submodel: submodel,
		section:  o.section,
		behind:   o.behind,
	})
}

func (c *Coffee) Input(opts ...InputOption) (Promise[string], error) {
	p, resolve := promise.New[string]()
	input := configuredInput(defaultInputOptions().apply(opts...)).withResolve(resolve)
	if err := c.AddSubmodel(input); err != nil {
		var zero Promise[string]
		return zero, err
	}
	return p, nil
}

func (c *Coffee) AwaitInput(opts ...InputOption) (string, error) {
	p, err := c.Input(opts...)
	if err != nil {
		return "", err
	}

	value, err := p.Await(c.ctx)
	return value, err
}

func (c *Coffee) ConfirmPromise(prompt string, confirm bool, opts ...SubmodelOption) (Promise[bool], error) {
	p, resolve := promise.New[bool]()
	input := NewConfirm(prompt, confirm).withResolve(resolve)
	if err := c.AddSubmodel(input, opts...); err != nil {
		var zero Promise[bool]
		return zero, err
	}
	return p, nil
}

func (c *Coffee) Confirm(prompt string, confirm bool, opts ...SubmodelOption) (bool, error) {
	p, err := c.ConfirmPromise(prompt, confirm, opts...)
	if err != nil {
		return false, err
	}

	value, err := p.Await(c.ctx)
	return value, err
}

func (c *Coffee) ConfirmAsync(prompt string, opts ...SubmodelOption) (Promise[bool], error) {
	return c.ConfirmPromise(prompt, true, opts...)
}

func (c *Coffee) AwaitConfirm(prompt string, opts ...SubmodelOption) (bool, error) {
	return c.Confirm(prompt, true, opts...)
}

func (c *Coffee) Select(prompt string, options []string) (Promise[string], error) {
	if len(options) == 0 {
		return Promise[string]{}, fmt.Errorf("select options cannot be empty")
	}
	return c.selectAt(prompt, options, 0)
}

func (c *Coffee) AwaitSelect(prompt string, options []string) (string, error) {
	p, err := c.Select(prompt, options)
	if err != nil {
		return "", err
	}

	value, err := p.Await(c.ctx)
	return value, err
}

func (c *Coffee) SelectDefault(prompt string, options []string, defaultValue string) (Promise[string], error) {
	if len(options) == 0 {
		return Promise[string]{}, fmt.Errorf("select options cannot be empty")
	}

	index := slices.Index(options, defaultValue)
	if index == -1 {
		return Promise[string]{}, fmt.Errorf("default value not found in options")
	}

	return c.selectAt(prompt, options, index)
}

func (c *Coffee) selectAt(prompt string, options []string, index int) (Promise[string], error) {
	p, resolve := promise.New[string]()
	selectModel := NewSelect(prompt, slices.Clone(options), index).withResolve(resolve)
	if err := c.AddSubmodel(selectModel); err != nil {
		var zero Promise[string]
		return zero, err
	}

	return p, nil
}

func (c *Coffee) AwaitSelectDefault(prompt string, options []string, defaultValue string) (string, error) {
	p, err := c.SelectDefault(prompt, options, defaultValue)
	if err != nil {
		return "", err
	}

	value, err := p.Await(c.ctx)
	return value, err
}

func (c *Coffee) send(msg any) error {
	if err := causeOrErr(c.ctx); err != nil {
		return err
	}

	c.program.Send(msg)
	return nil
}
