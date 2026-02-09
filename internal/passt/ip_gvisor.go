//go:build gvisor

package passt

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/tcpip"
	"gvisor.dev/gvisor/pkg/tcpip/network/arp"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv4"
	"gvisor.dev/gvisor/pkg/tcpip/network/ipv6"
	"gvisor.dev/gvisor/pkg/tcpip/stack"
	"gvisor.dev/gvisor/pkg/tcpip/transport/icmp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/tcp"
	"gvisor.dev/gvisor/pkg/tcpip/transport/udp"
)

// gvisorStack keeps the heavy dependency optional while preserving migration path.
type gvisorStack struct {
	*stack.Stack
}

func (n *Netstack) InitGVisor() *gvisorStack {
	s := stack.New(stack.Options{
		NetworkProtocols: []stack.NetworkProtocolFactory{ipv4.NewProtocol, ipv6.NewProtocol, arp.NewProtocol},
		TransportProtocols: []stack.TransportProtocolFactory{
			tcp.NewProtocol,
			udp.NewProtocol,
			icmp.NewProtocol4,
			icmp.NewProtocol6,
		},
	})
	return &gvisorStack{Stack: s}
}

func (g *gvisorStack) AttachNIC(id tcpip.NICID, ep stack.LinkEndpoint) error {
	if err := g.Stack.CreateNIC(id, ep); err != nil {
		return fmt.Errorf("create NIC %d: %w", id, err)
	}
	return nil
}
