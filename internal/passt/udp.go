package passt

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// UDPProxy ports the datagram forwarding role of udp.c.
type UDPProxy struct {
	ListenAddr string
	TargetAddr string
	conn       *net.UDPConn
	targetConn *net.UDPConn
	targetUDP  *net.UDPAddr
	wg         sync.WaitGroup
	log        *Logger

	mu      sync.Mutex
	clients map[string]*net.UDPAddr
}

func NewUDPProxy(listenAddr, targetAddr string, logger *Logger) *UDPProxy {
	return &UDPProxy{ListenAddr: listenAddr, TargetAddr: targetAddr, log: logger, clients: map[string]*net.UDPAddr{}}
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
	p.targetUDP, err = net.ResolveUDPAddr("udp", p.TargetAddr)
	if err != nil {
		return fmt.Errorf("resolve udp target addr: %w", err)
	}
	p.targetConn, err = net.DialUDP("udp", nil, p.targetUDP)
	if err != nil {
		return fmt.Errorf("dial udp target %s: %w", p.TargetAddr, err)
	}

	p.wg.Add(2)
	go p.forwardClientToTarget(ctx)
	go p.forwardTargetToClient(ctx)
	return nil
}

func (p *UDPProxy) forwardClientToTarget(ctx context.Context) {
	defer p.wg.Done()
	buf := make([]byte, 65535)
	for {
		n, client, err := p.conn.ReadFromUDP(buf)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			p.log.Errorf("udp read failed: %v", err)
			return
		}
		p.mu.Lock()
		p.clients[p.targetUDP.String()] = client
		p.mu.Unlock()
		if _, err := p.targetConn.Write(buf[:n]); err != nil {
			p.log.Errorf("udp write to target failed: %v", err)
		}
	}
}

func (p *UDPProxy) forwardTargetToClient(ctx context.Context) {
	defer p.wg.Done()
	buf := make([]byte, 65535)
	for {
		_ = p.targetConn.SetReadDeadline(time.Now().Add(300 * time.Millisecond))
		n, err := p.targetConn.Read(buf)
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if ctx.Err() != nil {
					return
				}
				continue
			}
			if ctx.Err() != nil {
				return
			}
			p.log.Errorf("udp read from target failed: %v", err)
			return
		}
		p.mu.Lock()
		client := p.clients[p.targetUDP.String()]
		p.mu.Unlock()
		if client == nil {
			continue
		}
		if _, err := p.conn.WriteToUDP(buf[:n], client); err != nil {
			p.log.Errorf("udp write to client failed: %v", err)
		}
	}
}

func (p *UDPProxy) Stop(context.Context) error {
	if p.conn != nil {
		_ = p.conn.Close()
	}
	if p.targetConn != nil {
		_ = p.targetConn.Close()
	}
	p.wg.Wait()
	return nil
}
