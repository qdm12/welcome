package docker

import "welcome/pkg/terminal"

func IsDockerInstalled() bool {
	_, err := terminal.RunCommand("docker")
	if err != nil {
		return false
	}
	return true
}

func IsDockerRunning() bool {
	_, err := terminal.RunCommand("docker", "ps")
	if err != nil {
		return false
	}
	return true
}

func IsDockerComposeInstalled() bool {
	_, err := terminal.RunCommand("docker-compose", "version")
	if err != nil {
		return false
	}
	return true
}
