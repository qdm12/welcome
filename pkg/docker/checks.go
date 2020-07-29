package docker

import (
	"context"
)

func (d *docker) IsInstalled(ctx context.Context) (installed bool) {
	_, err := d.commander.Run(ctx, "docker")
	return err == nil
}

func (d *docker) IsRunning(ctx context.Context) (running bool) {
	_, err := d.commander.Run(ctx, "docker", "ps")
	return err == nil
}

func (d *docker) IsComposeInstalled(ctx context.Context) (installed bool) {
	_, err := d.commander.Run(ctx, "docker-compose", "version")
	return err == nil
}
