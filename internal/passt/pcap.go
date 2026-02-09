package passt

import (
	"encoding/binary"
	"os"
	"sync"
	"time"
)

type PCAPWriter struct {
	mu sync.Mutex
	f  *os.File
}

func NewPCAPWriter(path string) (*PCAPWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	w := &PCAPWriter{f: f}
	hdr := make([]byte, 24)
	binary.LittleEndian.PutUint32(hdr[0:4], 0xa1b2c3d4)
	binary.LittleEndian.PutUint16(hdr[4:6], 2)
	binary.LittleEndian.PutUint16(hdr[6:8], 4)
	binary.LittleEndian.PutUint32(hdr[16:20], 65535)
	binary.LittleEndian.PutUint32(hdr[20:24], 1)
	if _, err := f.Write(hdr); err != nil {
		_ = f.Close()
		return nil, err
	}
	return w, nil
}

func (w *PCAPWriter) WritePacket(pkt []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	h := make([]byte, 16)
	binary.LittleEndian.PutUint32(h[0:4], uint32(now.Unix()))
	binary.LittleEndian.PutUint32(h[4:8], uint32(now.Nanosecond()/1000))
	binary.LittleEndian.PutUint32(h[8:12], uint32(len(pkt)))
	binary.LittleEndian.PutUint32(h[12:16], uint32(len(pkt)))
	if _, err := w.f.Write(h); err != nil {
		return err
	}
	_, err := w.f.Write(pkt)
	return err
}

func (w *PCAPWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.f == nil {
		return nil
	}
	err := w.f.Close()
	w.f = nil
	return err
}
