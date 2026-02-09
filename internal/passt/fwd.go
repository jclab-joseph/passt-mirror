package passt

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// L2Forwarder bridges two TAP-like endpoints using an FDB-style MAC table.
type L2Forwarder struct {
	guest TapDevice
	host  TapDevice
	fdb   *MACTable
	pcap  *PCAPWriter
	log   *Logger
	wg    sync.WaitGroup
}

func NewL2Forwarder(guest, host TapDevice, fdb *MACTable, pcap *PCAPWriter, logger *Logger) *L2Forwarder {
	return &L2Forwarder{guest: guest, host: host, fdb: fdb, pcap: pcap, log: logger}
}

func (f *L2Forwarder) Name() string { return "l2-forwarder" }

func (f *L2Forwarder) Start(ctx context.Context) error {
	if f.guest == nil || f.host == nil {
		return errors.New("tap endpoints are required")
	}
	if f.fdb == nil {
		f.fdb = NewMACTable(0)
	}
	f.wg.Add(2)
	go f.pump(ctx, f.guest, f.host)
	go f.pump(ctx, f.host, f.guest)
	return nil
}

func (f *L2Forwarder) pump(ctx context.Context, src, dst TapDevice) {
	defer f.wg.Done()
	for {
		frameBytes, err := src.ReadFrame(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			f.log.Errorf("tap read %s failed: %v", src.Name(), err)
			return
		}
		frame, err := ParseEthernetFrame(frameBytes)
		if err != nil {
			f.log.Errorf("drop malformed frame on %s: %v", src.Name(), err)
			continue
		}
		f.fdb.Learn(frame.SrcMAC, src.Name())
		if f.pcap != nil {
			_ = f.pcap.WritePacket(frameBytes)
		}
		if err := dst.WriteFrame(ctx, frameBytes); err != nil {
			if ctx.Err() != nil {
				return
			}
			f.log.Errorf("tap write %s failed: %v", dst.Name(), err)
			return
		}
	}
}

func (f *L2Forwarder) Stop(context.Context) error {
	f.wg.Wait()
	if f.pcap != nil {
		if err := f.pcap.Close(); err != nil {
			return fmt.Errorf("close pcap: %w", err)
		}
	}
	return nil
}
