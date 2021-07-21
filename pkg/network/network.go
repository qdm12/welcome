package network

import (
	"context"
	"net"

	"github.com/qdm12/ddns-updater/pkg/publicip"
	"github.com/qdm12/golibs/connectivity"
)

var _ NetworkInterface = new(Network)

type NetworkInterface interface {
	OutboundIP(ctx context.Context) (ip string, err error)
	PublicIP(ctx context.Context) (ip string, err error)
	Check(ctx context.Context, urls []string) (errors []error)
}

type Network struct {
	pubipFetcher publicip.Fetcher
	connChecker  connectivity.Checker
}

func New(resolver *net.Resolver) (n *Network, err error) {
	pubipFetcher, err := publicip.NewFetcher(
		publicip.DNSSettings{Enabled: true}, publicip.HTTPSettings{})
	if err != nil {
		return nil, err
	}
	return &Network{
		pubipFetcher: pubipFetcher,
		connChecker:  connectivity.NewDNSChecker(resolver),
	}, nil
}
