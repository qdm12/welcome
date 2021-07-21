package network

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

var ErrPrivateIP = errors.New("cannot get private IP address")

// Get preferred outbound ip of this machine
// TODO find from routing table
func (n *Network) OutboundIP(ctx context.Context) (ip string, err error) {
	d := &net.Dialer{Timeout: time.Second}
	conn, err := d.DialContext(ctx, "udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrPrivateIP, err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func (n *Network) PublicIP(ctx context.Context) (ip string, err error) {
	netIP, err := n.pubipFetcher.IP(ctx)
	if err != nil {
		return "", err
	}
	return netIP.String(), nil
}
