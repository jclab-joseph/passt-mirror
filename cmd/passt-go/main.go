package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/containers/passt-go/internal/passt"
)

func main() {
	cfg := passt.ParseConfig()
	engine := passt.NewEngine(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := engine.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}
}
