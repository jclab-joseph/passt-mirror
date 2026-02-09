package passt

import "context"

// Runtime collects runtime resources similar to global state in the C codebase.
type Runtime struct {
	Config   Config
	Logger   *Logger
	Netstack *Netstack
}

// Service defines the lifecycle of major subsystems.
type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
}
