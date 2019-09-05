package docker

import "welcome/pkg/terminal"

func GetDockerVersion() string {
	version, err := terminal.RunCommand("docker", "version", "--format", "'{{.Server.Version}}'")
	if err != nil {
		return "N/A"
	}
	return version
}

func GetDockerComposeVersion() string {
	version, err := terminal.RunCommand("docker-compose", "version", "--short")
	if err != nil {
		return "N/A"
	}
	return version
}
