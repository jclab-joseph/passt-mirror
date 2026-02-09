package passt

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
)

// TCPProxy ports the responsibility of tcp.c at a high level:
// accept host-side TCP connections and forward to a configured target.
type TCPProxy struct {
	ListenAddr string
	TargetAddr string
	ln         net.Listener
	wg         sync.WaitGroup
	log        *Logger
}

func NewTCPProxy(listenAddr, targetAddr string, logger *Logger) *TCPProxy {
	return &TCPProxy{ListenAddr: listenAddr, TargetAddr: targetAddr, log: logger}
}

func (p *TCPProxy) Name() string { return "tcp" }

func (p *TCPProxy) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", p.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen tcp %s: %w", p.ListenAddr, err)
	}
	p.ln = ln
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				p.log.Errorf("tcp accept failed: %v", err)
				return
			}
			p.wg.Add(1)
			go func(c net.Conn) {
				defer p.wg.Done()
				p.handle(c)
			}(conn)
		}
	}()
	return nil
}

func (p *TCPProxy) handle(src net.Conn) {
	defer src.Close()
	dst, err := net.Dial("tcp", p.TargetAddr)
	if err != nil {
		p.log.Errorf("tcp dial %s failed: %v", p.TargetAddr, err)
		return
	}
	defer dst.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); _, _ = io.Copy(dst, src) }()
	go func() { defer wg.Done(); _, _ = io.Copy(src, dst) }()
	wg.Wait()
}

func (p *TCPProxy) Stop(context.Context) error {
	if p.ln != nil {
		_ = p.ln.Close()
	}
	p.wg.Wait()
	return nil
}
