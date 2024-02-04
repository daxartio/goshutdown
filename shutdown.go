package goshutdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrShutdown = errors.New("shutdown error")

const (
	DefaultTimeout = 10 * time.Second
)

type (
	Handler       = func(ctx context.Context) error
	NotifyContext = func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc)
)

type Shutdown struct {
	Timeout       time.Duration
	Handler       Handler
	Signals       []os.Signal
	NotifyContext NotifyContext
}

// New creates a new Shutdown instance.
//
// The default timeout is 10 seconds.
func New() *Shutdown {
	return &Shutdown{
		Timeout:       DefaultTimeout,
		Handler:       nil,
		Signals:       []os.Signal{syscall.SIGINT, syscall.SIGTERM},
		NotifyContext: signal.NotifyContext,
	}
}

// WithTimeout sets the timeout for the shutdown process.
func (s *Shutdown) WithTimeout(timeout time.Duration) *Shutdown {
	s.Timeout = timeout

	return s
}

// WithHandler sets the handler for the shutdown process.
func (s *Shutdown) WithHandler(handler Handler) *Shutdown {
	s.Handler = handler

	return s
}

// WithNotifyContext sets the notify context function.
func (s *Shutdown) WithNotifyContext(notifyContext NotifyContext) *Shutdown {
	s.NotifyContext = notifyContext

	return s
}

// WithSignals sets the signals to listen for.
func (s *Shutdown) WithSignals(signals ...os.Signal) *Shutdown {
	s.Signals = signals

	return s
}

// Wait waits for the shutdown signal to be received or the timeout to expire.
// It returns an error if the shutdown process encounters an error or if the timeout is exceeded.
func (s *Shutdown) Wait() error {
	ctx, stop := s.NotifyContext(context.Background(), s.Signals...)
	defer stop()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	done := make(chan error)

	if s.Handler != nil {
		go func() {
			defer close(done)

			err := s.Handler(ctx)
			if err != nil {
				done <- err
			}
		}()
	}

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("%w: %w", ErrShutdown, err)
		}
	case <-ctx.Done():
	}

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("%w: %w", ErrShutdown, err)
	}

	return nil
}
