package passt

import (
	"context"
	"errors"
	"time"
)

// Engine is the Go equivalent of passt.c + event-loop centric modules.
type Engine struct {
	runtime Runtime
}

func NewEngine(cfg Config) *Engine {
	logger := NewLogger()
	ns := NewNetstack(cfg)
	return &Engine{runtime: Runtime{Config: cfg, Logger: logger, Netstack: ns}}
}

func (e *Engine) Run(ctx context.Context) error {
	e.runtime.Logger.Infof("starting mode=%s addr=%s mtu=%d", e.runtime.Config.Mode, e.runtime.Config.ListenAddr, e.runtime.Config.MTU)
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), e.runtime.Config.ShutdownTimout)
	defer cancel()
	select {
	case <-shutdownCtx.Done():
	}
	if errors.Is(ctx.Err(), context.Canceled) {
		return context.Canceled
	}
	return nil
}

func (e *Engine) Health() map[string]string {
	return map[string]string{
		"mode":   e.runtime.Config.Mode,
		"uptime": time.Now().Format(time.RFC3339),
	}
}
