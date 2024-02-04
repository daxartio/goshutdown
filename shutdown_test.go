package goshutdown_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/daxartio/goshutdown"
)

var ErrTest = errors.New("test error")

func TestShutdown(t *testing.T) {
	t.Parallel()

	err := goshutdown.New().
		WithTimeout(goshutdown.DefaultTimeout).
		WithHandler(func(ctx context.Context) error {
			return nil
		}).
		WithNotifyContext(func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			return ctx, func() {}
		}).
		Wait()
	if err != nil {
		t.Error(err)
	}
}

func TestShutdownError(t *testing.T) {
	t.Parallel()

	err := goshutdown.New().
		WithTimeout(goshutdown.DefaultTimeout).
		WithHandler(func(ctx context.Context) error {
			return ErrTest
		}).
		WithNotifyContext(func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			return ctx, func() {}
		}).
		Wait()
	if err == nil {
		t.Error("expected error")
	}

	if !errors.Is(err, ErrTest) {
		t.Error("expected test error")
	}

	if !errors.Is(err, goshutdown.ErrShutdown) {
		t.Error("expected shutdown error")
	}
}

func TestShutdownTimeout(t *testing.T) {
	t.Parallel()

	err := goshutdown.New().
		WithTimeout(time.Microsecond).
		WithHandler(func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)

			return nil
		}).
		WithNotifyContext(func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			return ctx, func() {}
		}).
		Wait()
	if err == nil {
		t.Error("expected error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Error("expected deadline exceeded error")
	}

	if !errors.Is(err, goshutdown.ErrShutdown) {
		t.Error("expected shutdown error")
	}
}
