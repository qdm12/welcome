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
	commander command.Commander
	client    client.APIClient
}

func New(commander command.Commander) (d Docker, err error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, fmt.Errorf("cannot use Docker: %w", err)
	}
	return &docker{
		commander: commander,
		client:    dockerClient,
	}, nil
}
