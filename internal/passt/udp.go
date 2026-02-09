package passt

import (
	"context"
	"fmt"
	"net"
	"sync"
)

// UDPProxy ports the datagram forwarding role of udp.c.
type UDPProxy struct {
	ListenAddr string
	TargetAddr string
	conn       *net.UDPConn
	wg         sync.WaitGroup
	log        *Logger
}

func NewUDPProxy(listenAddr, targetAddr string, logger *Logger) *UDPProxy {
	return &UDPProxy{ListenAddr: listenAddr, TargetAddr: targetAddr, log: logger}
}

func (p *UDPProxy) Name() string { return "udp" }

func (p *UDPProxy) Start(ctx context.Context) error {
	addr, err := net.ResolveUDPAddr("udp", p.ListenAddr)
	if err != nil {
		return fmt.Errorf("resolve udp listen addr: %w", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("listen udp %s: %w", p.ListenAddr, err)
	}
	p.conn = conn
	target, err := net.ResolveUDPAddr("udp", p.TargetAddr)
	if err != nil {
		return fmt.Errorf("resolve udp target addr: %w", err)
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		buf := make([]byte, 65535)
		for {
			n, client, err := conn.ReadFromUDP(buf)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				p.log.Errorf("udp read failed: %v", err)
				return
			}
			if _, err := conn.WriteToUDP(buf[:n], target); err != nil {
				p.log.Errorf("udp write to target failed: %v", err)
				continue
			}
			_, _ = conn.WriteToUDP([]byte("ok"), client)
		}
	}()
	return nil
}

func (p *UDPProxy) Stop(context.Context) error {
	if p.conn != nil {
		_ = p.conn.Close()
	}
	p.wg.Wait()
	return nil
}
