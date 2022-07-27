package docker

import (
	"context"
	"os/exec"
)

func (d *Docker) Version(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "docker", "version", "--format", "'{{.Server.Version}}'")
	version, err := d.commander.Run(cmd)
	if err != nil {
		return "N/A"
	}
	return version
}

func (d *Docker) ComposeVersion(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "docker-compose", "version", "--short")
	version, err := d.commander.Run(cmd)
	if err != nil {
		return "N/A"
	}
	return version
}
