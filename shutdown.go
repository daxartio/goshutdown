package goshutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DefaultTimeout = 10 * time.Second
)

type Handler = func(ctx context.Context)

type Shutdown struct {
	Timeout time.Duration
	Handler Handler
	Signals []os.Signal
}

// New creates a new Shutdown instance.
//
// The default timeout is 10 seconds.
func New() *Shutdown {
	return &Shutdown{
		Timeout: DefaultTimeout,
		Handler: nil,
		Signals: []os.Signal{syscall.SIGINT, syscall.SIGTERM},
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

// WithSignals sets the signals to listen for.
func (s *Shutdown) WithSignals(signals ...os.Signal) *Shutdown {
	s.Signals = signals

	return s
}

// Wait waits for a signal and calls the handler.
func (s *Shutdown) Wait() {
	ctx, stop := signal.NotifyContext(context.Background(), s.Signals...)
	defer stop()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	if s.Handler != nil {
		s.Handler(ctx)
	}
}
