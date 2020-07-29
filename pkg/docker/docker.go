package docker

import (
	"context"

	"github.com/qdm12/golibs/command"
)

type Docker interface {
	IsInstalled(ctx context.Context) (installed bool)
	IsRunning(ctx context.Context) (running bool)
	IsComposeInstalled(ctx context.Context) (installed bool)
	Version(ctx context.Context) string
	ComposeVersion(ctx context.Context) string
	CountContainers(ctx context.Context) (count int, err error)
	AreContainerRunning(ctx context.Context, requiredContainerNames []string) (containersNotRunning []string, err error)
	BadContainers(ctx context.Context) (badStatus []string, err error)
}

type docker struct {
	commander command.Commander
}

func New(commander command.Commander) Docker {
	return &docker{
		commander: commander,
	}
}
