# Graceful shutdown

[![Go Documentation](https://godocs.io/github.com/daxartio/goshutdown?status.svg)](https://godocs.io/github.com/daxartio/goshutdown)

The graceful shutdown is a mechanism that allows the application to finish the
current requests before shutting down. This is useful to avoid losing data or
to avoid corrupting the data.

The package provides an abstraction to handle the graceful shutdown.
It listens to the `SIGINT` and `SIGTERM` signals and calls the `Shutdown` method.

## Usage

```go
package main

import (
	"context"
	"time"

	"github.com/daxartio/goshutdown"
)

func main() {
	server := &Server{}

	go server.Run()

	err := goshutdown.New().
		WithTimeout(10 * time.Second).
		WithHandler(func(ctx context.Context) error {
			println("Shutting down...")

			server.Stop(ctx)

			return nil
		}).Wait()
	if err != nil {
		println(err.Error())
	}
}
```
