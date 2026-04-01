package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/olimci/coffee"
)

var demoBindings = []coffee.Keybind{
	{Key: "i", Event: "input", Description: "input"},
	{Key: "d", Event: "discard", Description: "discard input"},
	{Key: "y", Event: "confirm", Description: "confirm clear"},
	{Key: "m", Event: "choice", Description: "mode select"},
	{Key: "r", Event: "rebuild", Description: "rebuild"},
	{Key: "f", Event: "fail", Description: "fail build"},
	{Key: "l", Event: "log", Description: "log line"},
	{Key: "c", Event: "clear", Description: "clear logs"},
	{Key: "q", Event: "quit", Description: "quit"},
}

func main() {
	ctx := context.Background()

	err := coffee.Do(func(ctx context.Context, c *coffee.Coffee) error {
		app := &demoApp{}

		_ = c.SetWindowTitle("coffee test-app")
		_ = c.Log("test-app: dev shell demo")
		_ = c.Log("Body-less version of shizuka's dev surface.")
		_ = c.Log("Use the footer keys to trigger input, confirm, choice, loading, logs, and clearing.")

		return app.run(ctx, c)
	}, coffee.WithContext(ctx))
	if err == nil {
		return
	}

	if errors.Is(err, coffee.ErrNonInteractive) {
		fmt.Fprintln(os.Stderr, "test-app requires an interactive terminal")
		return
	}

	panic(fmt.Errorf("test-app failed: %w", err))
}

type demoApp struct {
	builds int
	logs   int
}

func (a *demoApp) run(ctx context.Context, c *coffee.Coffee) error {
	status, err := c.Status("watching for changes")
	if err != nil {
		return err
	}

	keys, err := c.Keybinds(
		demoBindings,
	)
	if err != nil {
		return err
	}
	defer func() {
		_ = keys.Clear()
		_ = status.Clear()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-keys.Events():
			if !ok {
				return nil
			}

			switch ev.Event {
			case "input":
				if err := a.captureInput(ctx, c, status); err != nil {
					return err
				}
			case "discard":
				if err := a.captureDiscardInput(c); err != nil {
					return err
				}
			case "confirm":
				if err := a.confirmClear(c, status); err != nil {
					return err
				}
			case "choice":
				if err := a.chooseAction(ctx, c, status); err != nil {
					return err
				}
			case "rebuild":
				if err := a.simulateBuild(ctx, c, status, "manual", false); err != nil {
					return err
				}
			case "fail":
				if err := a.simulateBuild(ctx, c, status, "manual", true); err != nil {
					return err
				}
			case "log":
				a.logs++
				_ = c.Logf("[%02d] synthetic log line from demo action", a.logs)
				_ = status.Idle("watching for changes")
			case "clear":
				_ = c.Clear()
				_ = status.Idle("watching for changes")
			case "quit":
				_ = status.Success("stopped")
				return nil
			}
		}
	}
}

func (a *demoApp) captureInput(ctx context.Context, c *coffee.Coffee, status *coffee.StatusHandle) error {
	input, err := c.AwaitInput(
		coffee.WithInputPlaceholder("rebuild, fail, clear, or anything else"),
		coffee.WithInputSuggestions([]string{
			"rebuild",
			"fail",
			"clear",
			"touch content/index.md",
		}),
		coffee.WithInputWidth(36),
	)
	if err != nil {
		return err
	}

	command := strings.TrimSpace(input)
	switch command {
	case "":
		_ = c.Log("input submitted: <empty>")
	case "clear":
		_ = c.Clear()
		_ = status.Idle("watching for changes")
	case "rebuild":
		_ = c.Log(`input submitted: "rebuild"`)
		if err := a.simulateBuild(ctx, c, status, "input", false); err != nil {
			return err
		}
	case "fail":
		_ = c.Log(`input submitted: "fail"`)
		if err := a.simulateBuild(ctx, c, status, "input", true); err != nil {
			return err
		}
	default:
		_ = c.Logf("input submitted: %q", command)
	}

	return nil
}

func (a *demoApp) captureDiscardInput(c *coffee.Coffee) error {
	input, err := c.AwaitInput(
		coffee.WithInputPlaceholder("type a note that should not stay on screen"),
		coffee.WithInputWidth(36),
		coffee.WithInputDiscardSubmitted(),
	)
	if err != nil {
		return err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return c.Log("discard input submitted: <empty>")
	}

	return c.Logf("discard input submitted: %q", input)
}

func (a *demoApp) confirmClear(c *coffee.Coffee, status *coffee.StatusHandle) error {
	confirmed, err := c.Confirm("clear logs?", false)
	if err != nil {
		return err
	}

	if !confirmed {
		return c.Log("clear cancelled")
	}

	if err := c.Clear(); err != nil {
		return err
	}

	return status.Idle("watching for changes")
}

func (a *demoApp) chooseAction(ctx context.Context, c *coffee.Coffee, status *coffee.StatusHandle) error {
	choice, err := c.AwaitSelectDefault("choose action", []string{
		"watch",
		"rebuild",
		"fail build",
		"clear logs",
	}, "watch")
	if err != nil {
		return err
	}

	switch choice {
	case "watch":
		if err := status.Idle("watching for changes"); err != nil {
			return err
		}
		return c.Log("mode select: watch")
	case "rebuild":
		_ = c.Log(`mode select: "rebuild"`)
		return a.simulateBuild(ctx, c, status, "choice", false)
	case "fail build":
		_ = c.Log(`mode select: "fail build"`)
		return a.simulateBuild(ctx, c, status, "choice", true)
	case "clear logs":
		if err := c.Clear(); err != nil {
			return err
		}
		if err := status.Idle("watching for changes"); err != nil {
			return err
		}
		return c.Log("mode select: clear logs")
	default:
		return c.Logf("mode select: %q", choice)
	}
}

func (a *demoApp) simulateBuild(ctx context.Context, c *coffee.Coffee, status *coffee.StatusHandle, trigger string, fail bool) error {
	a.builds++
	start := time.Now()

	if err := c.Clear(); err != nil {
		return err
	}
	if err := c.Logf("build #%d triggered by %s", a.builds, trigger); err != nil {
		return err
	}

	if err := status.Working(fmt.Sprintf("building (%s)", trigger)); err != nil {
		return err
	}
	if err := sleepContext(ctx, 300*time.Millisecond); err != nil {
		return err
	}

	if err := status.Progress("loading config", 0.0); err != nil {
		return err
	}

	steps := []struct {
		percent float64
		status  string
		logLine string
		delay   time.Duration
	}{
		{percent: 0.2, status: "loading config", logLine: "loaded config", delay: 180 * time.Millisecond},
		{percent: 0.45, status: "scanning content", logLine: "scanned content tree", delay: 220 * time.Millisecond},
		{percent: 0.7, status: "rendering templates", logLine: "rendered templates", delay: 220 * time.Millisecond},
		{percent: 1.0, status: "writing output", logLine: "wrote output files", delay: 180 * time.Millisecond},
	}

	for _, step := range steps {
		if err := sleepContext(ctx, step.delay); err != nil {
			return err
		}
		if err := status.Message(step.status); err != nil {
			return err
		}
		if err := status.SetProgress(step.percent); err != nil {
			return err
		}
		_ = c.Log(step.logLine)
	}

	elapsed := time.Since(start).Truncate(time.Millisecond)
	if fail {
		if err := status.Error(fmt.Sprintf("build failed (%s)", elapsed)); err != nil {
			return err
		}
		_ = c.Logf("build failed (%s)", elapsed)
		_ = c.Log("demo error: template render failed near content/posts/demo.md")
		return nil
	}

	if err := status.Success(fmt.Sprintf("build complete (%s)", elapsed)); err != nil {
		return err
	}
	_ = c.Logf("build complete (%s)", elapsed)
	_ = c.Log("reload event broadcast")

	if err := sleepContext(ctx, 400*time.Millisecond); err != nil {
		return err
	}

	return status.Idle("watching for changes")
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
