package coffee

import (
	"context"
	"errors"
	"os"
	"strings"
)

func isNonInteractive(err error) bool {
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		switch strings.ToLower(pathErr.Path) {
		case "/dev/tty", "conin$":
			return true
		}
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "could not open a new tty") ||
		strings.Contains(msg, "open /dev/tty") ||
		strings.Contains(msg, "open conin$")
}

func causeOrErr(ctx context.Context) error {
	if cause := context.Cause(ctx); cause != nil {
		return cause
	}
	return ctx.Err()
}

func cleanCancel(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, context.Canceled) {
		if cause := context.Cause(ctx); cause != nil {
			if errors.Is(cause, context.Canceled) {
				return nil
			}
			return cause
		}
	}

	return err
}
