package passt

import (
	"context"
	"errors"
	"sync"
)

// TapDevice models an L2 endpoint similar to Linux TAP semantics.
type TapDevice interface {
	ReadFrame(ctx context.Context) ([]byte, error)
	WriteFrame(ctx context.Context, frame []byte) error
	Close() error
	Name() string
}

// MemoryTap is a cross-platform TAP-like endpoint used for tests and user-space wiring.
type MemoryTap struct {
	name   string
	in     chan []byte
	peer   *MemoryTap
	mu     sync.RWMutex
	closed chan struct{}
}

func NewMemoryTapPair(a, b string, depth int) (*MemoryTap, *MemoryTap) {
	if depth < 1 {
		depth = 64
	}
	ta := &MemoryTap{name: a, in: make(chan []byte, depth), closed: make(chan struct{})}
	tb := &MemoryTap{name: b, in: make(chan []byte, depth), closed: make(chan struct{})}
	ta.peer = tb
	tb.peer = ta
	return ta, tb
}

func (t *MemoryTap) Name() string { return t.name }

func (t *MemoryTap) ReadFrame(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-t.closed:
		return nil, errors.New("tap closed")
	case f := <-t.in:
		return append([]byte(nil), f...), nil
	}
}

func (t *MemoryTap) WriteFrame(ctx context.Context, frame []byte) error {
	t.mu.RLock()
	peer := t.peer
	t.mu.RUnlock()
	if peer == nil {
		return errors.New("tap disconnected")
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.closed:
		return errors.New("tap closed")
	case peer.in <- append([]byte(nil), frame...):
		return nil
	}
}

func (t *MemoryTap) Close() error {
	select {
	case <-t.closed:
		return nil
	default:
		close(t.closed)
		return nil
	}
}
