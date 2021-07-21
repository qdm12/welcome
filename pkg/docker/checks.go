package docker

import (
	"context"
	"os/exec"
)

func (d *docker) IsRunning(ctx context.Context) (running bool) {
	_, err := d.client.Ping(ctx)
	return err == nil
}

func (d *docker) IsComposeInstalled(ctx context.Context) (installed bool) {
	cmd := exec.CommandContext(ctx, "docker-compose", "version")
	_, err := d.commander.Run(cmd)
	return err == nil
}
