package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/daxartio/goshutdown"
)

func main() {
	// Create a new Shutdown instance
	shutdown := goshutdown.New().WithTimeout(goshutdown.DefaultTimeout)

	ctx, cancel := context.WithCancel(context.Background())

	go run(ctx)

	// Set a shutdown handler
	shutdown.WithHandler(func(ctx context.Context) error {
		defer cancel()
		slog.InfoContext(ctx, "Shutting down...")

		return nil
	})

	slog.Info("Waiting for a signal...")

	// Wait for a signal
	err := shutdown.Wait()
	if err != nil {
		slog.Error(err.Error())
	}
}

func run(ctx context.Context) {
	slog.InfoContext(ctx, "Running...")

	const interval = 1 * time.Second

	for {
		slog.InfoContext(ctx, "Tick")
		time.Sleep(interval)
	}
}
