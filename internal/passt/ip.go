package passt

import (
	"encoding/binary"
	"fmt"
	"net"
)

import gchecksum "gvisor.dev/gvisor/pkg/tcpip/checksum"

// Netstack keeps gVisor TCP/IP primitives as checksum/parsing helpers.
type Netstack struct {
	cfg Config
}

func NewNetstack(cfg Config) *Netstack {
	return &Netstack{cfg: cfg}
}

func (n *Netstack) Validate() error {
	if net.ParseIP(n.cfg.ListenAddr) == nil {
		return fmt.Errorf("invalid listen address: %s", n.cfg.ListenAddr)
	}
	return nil
}

func IPv4HeaderChecksum(header []byte) uint16 {
	return uint16(gchecksum.Checksum(header, 0))
}

func PseudoHeaderChecksum(src, dst net.IP, proto uint8, length uint16) uint16 {
	var pseudo [12]byte
	copy(pseudo[0:4], src.To4())
	copy(pseudo[4:8], dst.To4())
	pseudo[8] = 0
	pseudo[9] = proto
	binary.BigEndian.PutUint16(pseudo[10:12], length)
	return uint16(gchecksum.Checksum(pseudo[:], 0))
}
