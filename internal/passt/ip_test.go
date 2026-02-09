package passt

import (
	"net"
	"testing"
)

func TestIPv4HeaderChecksum(t *testing.T) {
	hdr := []byte{0x45, 0x00, 0x00, 0x54, 0xa6, 0xf2, 0x40, 0x00, 0x40, 0x01, 0x00, 0x00, 0xc0, 0xa8, 0x00, 0x01, 0xc0, 0xa8, 0x00, 0xc7}
	c := IPv4HeaderChecksum(hdr)
	if c == 0 {
		t.Fatal("checksum should not be zero")
	}
}

func TestPseudoHeaderChecksum(t *testing.T) {
	c := PseudoHeaderChecksum(net.IPv4(10, 0, 0, 1), net.IPv4(10, 0, 0, 2), 6, 20)
	if c == 0 {
		t.Fatal("pseudo checksum should not be zero")
	}
}
