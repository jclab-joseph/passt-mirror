package passt

// Netstack is the platform-neutral façade used by the engine.
// gVisor-backed implementation is provided in ip_gvisor.go under `-tags gvisor`.
type Netstack struct {
	cfg Config
}

func NewNetstack(cfg Config) *Netstack {
	return &Netstack{cfg: cfg}
}
