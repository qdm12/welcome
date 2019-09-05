package docker

import (
	"strings"
	"welcome/pkg/terminal"
	"welcome/pkg/utils"
)

func CountContainers() (count int, err error) {
	dockerPsNames, err := terminal.RunCommand("docker", "ps", "--format", "'{{.Names}}'")
	if err != nil {
		return 0, err
	}
	containersCount := utils.CountNonEmptyLines(dockerPsNames)
	return containersCount, nil
}

func IsContainerRunning(requiredContainerNames ...string) (containersNotRunning []string, err error) {
	dockerPsNames, err := terminal.RunCommand("docker", "ps", "--format", "'{{.Names}}'")
	if err != nil {
		return nil, err
	}
	containerNames := make(map[string]bool)
	for _, name := range utils.StringToLines(dockerPsNames) {
		containerNames[name] = true
	}
	for _, name := range requiredContainerNames {
		if _, ok := containerNames[name]; !ok {
			containersNotRunning = append(containersNotRunning, name)
		}
	}
	return containersNotRunning, nil
}

func GetBadContainers() (badStatus []string, err error) {
	dockerPsStatus, err := terminal.RunCommand("docker", "ps", "--format", "'Container {{.Names}} is {{.Status}}'")
	if err != nil {
		return nil, err
	}
	allStatus := utils.StringToLines(dockerPsStatus)
	for _, status := range allStatus {
		if strings.Contains(status, "unhealthy") || strings.Contains(status, "restarting") {
			badStatus = append(badStatus, status)
		}
	}
	return badStatus, nil
}
