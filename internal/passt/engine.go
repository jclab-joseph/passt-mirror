package passt

import (
	"context"
	"errors"
	"time"
)

// Engine is the Go equivalent of passt.c + event-loop centric modules.
type Engine struct {
	runtime  Runtime
	services []Service
}

func NewEngine(cfg Config) *Engine {
	logger := NewLogger()
	ns := NewNetstack(cfg)
	services := []Service{
		NewTCPProxy(cfg.TCPListen, cfg.TCPTarget, logger),
		NewUDPProxy(cfg.UDPListen, cfg.UDPTarget, logger),
	}
	return &Engine{runtime: Runtime{Config: cfg, Logger: logger, Netstack: ns}, services: services}
}

func (e *Engine) Run(ctx context.Context) error {
	e.runtime.Logger.Infof("starting mode=%s addr=%s mtu=%d", e.runtime.Config.Mode, e.runtime.Config.ListenAddr, e.runtime.Config.MTU)
	for _, svc := range e.services {
		if err := svc.Start(ctx); err != nil {
			return err
		}
	}
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), e.runtime.Config.ShutdownTimeout)
	defer cancel()
	for i := len(e.services) - 1; i >= 0; i-- {
		_ = e.services[i].Stop(shutdownCtx)
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
