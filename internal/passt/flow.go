package passt

import (
	"net"
	"sync"
	"time"
)

type MACTableEntry struct {
	Port     string
	LastSeen time.Time
}

type MACTable struct {
	mu      sync.RWMutex
	entries map[string]MACTableEntry
	age     time.Duration
}

func NewMACTable(age time.Duration) *MACTable {
	if age <= 0 {
		age = 5 * time.Minute
	}
	return &MACTable{entries: map[string]MACTableEntry{}, age: age}
}

func (t *MACTable) Learn(mac net.HardwareAddr, port string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries[mac.String()] = MACTableEntry{Port: port, LastSeen: time.Now()}
}

func (t *MACTable) Lookup(mac net.HardwareAddr) (string, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[mac.String()]
	if !ok || time.Since(e.LastSeen) > t.age {
		return "", false
	}
	return e.Port, true
}
