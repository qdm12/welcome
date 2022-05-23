package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/qdm12/golibs/command"
)

type Docker interface {
	IsRunning(ctx context.Context) (running bool)
	IsComposeInstalled(ctx context.Context) (installed bool)
	Version(ctx context.Context) string
	ComposeVersion(ctx context.Context) string
	CountContainers(ctx context.Context) (count int, err error)
	AreContainerRunning(ctx context.Context, requiredContainerNames []string) (containersNotRunning []string, err error)
	BadContainers(ctx context.Context) (containerNameToState map[string]string, err error)
}

type docker struct {
	commander Commander
	client    client.APIClient
}

type Commander interface {
	Run(cmd command.ExecCmd) (output string, err error)
}

func New(commander Commander) (
	d Docker, err error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		return nil, fmt.Errorf("cannot use Docker: %w", err)
	}
	return &docker{
		commander: commander,
		client:    dockerClient,
	}, nil
}
