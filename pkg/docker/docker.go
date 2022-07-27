package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/qdm12/golibs/command"
)

type Docker struct {
	commander Commander
	client    client.APIClient
}

type Commander interface {
	Run(cmd command.ExecCmd) (output string, err error)
}

func New(commander Commander) (
	d *Docker, err error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("cannot use Docker: %w", err)
	}
	return &Docker{
		commander: commander,
		client:    dockerClient,
	}, nil
}
