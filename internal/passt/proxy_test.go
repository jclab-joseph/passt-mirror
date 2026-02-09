package passt

import (
	"context"
	"io"
	"net"
	"testing"
	"time"
)

func TestTCPProxyForwards(t *testing.T) {
	targetLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer targetLn.Close()
	go func() {
		c, err := targetLn.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		buf := make([]byte, 4)
		if _, err := io.ReadFull(c, buf); err != nil {
			return
		}
		_, _ = c.Write(buf)
	}()

	proxy := NewTCPProxy("127.0.0.1:0", targetLn.Addr().String(), NewLogger())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := proxy.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer proxy.Stop(context.Background())

	conn, err := net.Dial("tcp", proxy.ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
	_, _ = conn.Write([]byte("ping"))
	buf := make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatal(err)
	}
	if string(buf) != "ping" {
		t.Fatalf("got %q", string(buf))
	}
}

func TestUDPProxyResponds(t *testing.T) {
	targetAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	targetConn, err := net.ListenUDP("udp", targetAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer targetConn.Close()

	proxy := NewUDPProxy("127.0.0.1:0", targetConn.LocalAddr().String(), NewLogger())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := proxy.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer proxy.Stop(context.Background())

	client, err := net.Dial("udp", proxy.conn.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	_ = client.SetDeadline(time.Now().Add(2 * time.Second))
	_, _ = client.Write([]byte("hello"))
	buf := make([]byte, 16)
	n, err := client.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "ok" {
		t.Fatalf("got %q", string(buf[:n]))
	}
}
