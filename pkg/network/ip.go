package network

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// Get preferred outbound ip of this machine
// TODO find from routing table
func (n *network) GetOutboundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("cannot get private IP address: %w", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

func (n *network) GetPublicIP() (ip string, err error) {
	content, status, err := n.client.GetContent("https://diagnostic.opendns.com/myip")
	if err != nil {
		return "", fmt.Errorf("cannot get public IP address: %w", err)
	} else if status != http.StatusOK {
		return "", fmt.Errorf("cannot get public IP address: HTTP status code %d", status)
	}
	ip = strings.TrimSpace(string(content))
	ip = strings.TrimPrefix(ip, "\n")
	ip = strings.TrimSuffix(ip, "\n")
	if ip == "" {
		return "", fmt.Errorf("Public IP address not found")
	}
	return ip, nil
}
