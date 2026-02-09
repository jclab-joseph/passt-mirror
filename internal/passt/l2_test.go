package passt

import (
	"context"
	"net"
	"os"
	"testing"
	"time"
)

func TestEthernetFrameRoundTrip(t *testing.T) {
	f := EthernetFrame{DstMAC: net.HardwareAddr{1, 2, 3, 4, 5, 6}, SrcMAC: net.HardwareAddr{6, 5, 4, 3, 2, 1}, EtherType: EtherTypeIPv4, Payload: []byte{0xaa, 0xbb}}
	b, err := f.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	p, err := ParseEthernetFrame(b)
	if err != nil {
		t.Fatal(err)
	}
	if p.EtherType != f.EtherType || p.DstMAC.String() != f.DstMAC.String() || p.SrcMAC.String() != f.SrcMAC.String() {
		t.Fatalf("frame mismatch: %#v %#v", p, f)
	}
}

func TestL2ForwarderForwardsFrames(t *testing.T) {
	guestSide, guestWire := NewMemoryTapPair("guest-side", "guest-wire", 8)
	hostWire, hostSide := NewMemoryTapPair("host-wire", "host-side", 8)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fwd := NewL2Forwarder(guestWire, hostWire, NewMACTable(0), nil, NewLogger())
	if err := fwd.Start(ctx); err != nil {
		t.Fatal(err)
	}
	defer func() {
		cancel()
		_ = fwd.Stop(context.Background())
	}()

	frame, _ := (EthernetFrame{DstMAC: net.HardwareAddr{0, 1, 2, 3, 4, 5}, SrcMAC: net.HardwareAddr{6, 7, 8, 9, 10, 11}, EtherType: EtherTypeARP, Payload: []byte{1, 2, 3}}).MarshalBinary()
	if err := guestSide.WriteFrame(ctx, frame); err != nil {
		t.Fatal(err)
	}
	readCtx, cancelRead := context.WithTimeout(ctx, time.Second)
	defer cancelRead()
	got, err := hostSide.ReadFrame(readCtx)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(frame) {
		t.Fatalf("unexpected frame length %d", len(got))
	}
}

func TestPCAPWriterWritesHeaderAndPacket(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "pcap-*.pcap")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	w, err := NewPCAPWriter(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if err := w.WritePacket([]byte{1, 2, 3, 4}); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	st, err := os.Stat(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if st.Size() <= 24 {
		t.Fatalf("pcap too small: %d", st.Size())
	}
}
