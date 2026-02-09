package passt

import (
	"flag"
	"time"
)

// Config mirrors the role of conf.c/conf.h.
type Config struct {
	Mode           string
	ListenAddr     string
	MTU            int
	EnablePCAP     bool
	PCAPPath       string
	ShutdownTimout time.Duration
}

func ParseConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.Mode, "mode", "passt", "operation mode: passt or pasta")
	flag.StringVar(&cfg.ListenAddr, "listen", "0.0.0.0", "listen address")
	flag.IntVar(&cfg.MTU, "mtu", 1500, "L2 MTU")
	flag.BoolVar(&cfg.EnablePCAP, "pcap", false, "enable packet capture output")
	flag.StringVar(&cfg.PCAPPath, "pcap-path", "passt.pcap", "pcap output path")
	flag.DurationVar(&cfg.ShutdownTimout, "shutdown-timeout", 5*time.Second, "graceful shutdown timeout")
	flag.Parse()
	return cfg
}
