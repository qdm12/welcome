package docker

import (
	"context"
	"os/exec"
	"strings"
)

func (d *Docker) Version(ctx context.Context) string {
	cmd := exec.CommandContext(ctx, "docker", "version", "--format", "'{{.Server.Version}}'")
	version, err := d.commander.Run(cmd)
	if err != nil {
		return "N/A"
	}
	return version
}

func (d *Docker) ComposeVersion(ctx context.Context) (version string) {
	composeCommands := []string{"docker compose", "docker-compose"}
	for _, composeCommand := range composeCommands {
		fields := strings.Fields(composeCommand)
		fields = append(fields, "version", "--short")
		cmd := exec.CommandContext(ctx, fields[0], fields[1:]...) //nolint:gosec
		version, err := d.commander.Run(cmd)
		if err == nil {
			return version
		}
	}

	return ""
}
