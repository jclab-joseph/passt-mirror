package passt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const EthernetHeaderLen = 14

type EtherType uint16

const (
	EtherTypeIPv4 EtherType = 0x0800
	EtherTypeARP  EtherType = 0x0806
	EtherTypeIPv6 EtherType = 0x86DD
)

type EthernetFrame struct {
	DstMAC    net.HardwareAddr
	SrcMAC    net.HardwareAddr
	EtherType EtherType
	Payload   []byte
}

func ParseEthernetFrame(b []byte) (EthernetFrame, error) {
	if len(b) < EthernetHeaderLen {
		return EthernetFrame{}, errors.New("short ethernet frame")
	}
	f := EthernetFrame{
		DstMAC:    append(net.HardwareAddr(nil), b[0:6]...),
		SrcMAC:    append(net.HardwareAddr(nil), b[6:12]...),
		EtherType: EtherType(binary.BigEndian.Uint16(b[12:14])),
		Payload:   append([]byte(nil), b[14:]...),
	}
	return f, nil
}

func (f EthernetFrame) MarshalBinary() ([]byte, error) {
	if len(f.DstMAC) != 6 || len(f.SrcMAC) != 6 {
		return nil, fmt.Errorf("invalid mac lens dst=%d src=%d", len(f.DstMAC), len(f.SrcMAC))
	}
	buf := make([]byte, EthernetHeaderLen+len(f.Payload))
	copy(buf[0:6], f.DstMAC)
	copy(buf[6:12], f.SrcMAC)
	binary.BigEndian.PutUint16(buf[12:14], uint16(f.EtherType))
	copy(buf[14:], f.Payload)
	return buf, nil
}
