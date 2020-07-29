package network

import (
	"time"

	netlib "github.com/qdm12/golibs/network"
)

type Network interface {
	GetOutboundIP() (ip string, err error)
	GetPublicIP() (ip string, err error)
	CheckMultipleHTTPConnections(urls []string) (errors []error)
}

type network struct {
	client netlib.Client
}

func New() Network {
	return &network{
		client: netlib.NewClient(time.Second),
	}
}
