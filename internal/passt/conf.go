package passt

import (
	"flag"
	"time"
)

// Config mirrors the role of conf.c/conf.h.
type Config struct {
	Mode            string
	ListenAddr      string
	MTU             int
	EnablePCAP      bool
	PCAPPath        string
	ShutdownTimeout time.Duration
	TCPListen       string
	TCPTarget       string
	UDPListen       string
	UDPTarget       string
}

func ParseConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.Mode, "mode", "passt", "operation mode: passt or pasta")
	flag.StringVar(&cfg.ListenAddr, "listen", "0.0.0.0", "listen address")
	flag.IntVar(&cfg.MTU, "mtu", 1500, "L2 MTU")
	flag.BoolVar(&cfg.EnablePCAP, "pcap", false, "enable packet capture output")
	flag.StringVar(&cfg.PCAPPath, "pcap-path", "passt.pcap", "pcap output path")
	flag.DurationVar(&cfg.ShutdownTimeout, "shutdown-timeout", 5*time.Second, "graceful shutdown timeout")
	flag.StringVar(&cfg.TCPListen, "tcp-listen", "127.0.0.1:10080", "TCP proxy listen address")
	flag.StringVar(&cfg.TCPTarget, "tcp-target", "127.0.0.1:80", "TCP proxy target address")
	flag.StringVar(&cfg.UDPListen, "udp-listen", "127.0.0.1:10053", "UDP proxy listen address")
	flag.StringVar(&cfg.UDPTarget, "udp-target", "127.0.0.1:53", "UDP proxy target address")
	flag.Parse()
	return cfg
}
