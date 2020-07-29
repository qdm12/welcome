package docker

import "context"

func (d *docker) Version(ctx context.Context) string {
	version, err := d.commander.Run(ctx, "docker", "version", "--format", "'{{.Server.Version}}'")
	if err != nil {
		return "N/A"
	}
	return version
}

func (d *docker) ComposeVersion(ctx context.Context) string {
	version, err := d.commander.Run(ctx, "docker-compose", "version", "--short")
	if err != nil {
		return "N/A"
	}
	return version
}
