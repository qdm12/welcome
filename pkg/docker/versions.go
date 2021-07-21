package docker

import (
	"context"
	"os/exec"
)

func (d *docker) Version(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "docker", "version", "--format", "'{{.Server.Version}}'")
	version, err := d.commander.Run(cmd)
	if err != nil {
		return "N/A"
	}
	return version
}

func (d *docker) ComposeVersion(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "docker-compose", "version", "--short")
	version, err := d.commander.Run(cmd)
	if err != nil {
		return "N/A"
	}
	return version
}
