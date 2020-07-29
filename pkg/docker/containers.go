package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/qdm12/welcome/pkg/utils"
)

func countNonEmptyLines(s string) (count int) {
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if len(line) > 0 {
			count++
		}
	}
	return count
}

func (d *docker) CountContainers(ctx context.Context) (count int, err error) {
	dockerPsNames, err := d.commander.Run(ctx, "docker", "ps", "--format", "'{{.Names}}'")
	if err != nil {
		return 0, fmt.Errorf("cannot list Docker containers: %w", err)
	}
	containersCount := countNonEmptyLines(dockerPsNames)
	return containersCount, nil
}

func (d *docker) AreContainerRunning(ctx context.Context, requiredContainerNames []string) (containersNotRunning []string, err error) {
	dockerPsNames, err := d.commander.Run(ctx, "docker", "ps", "--format", "'{{.Names}}'")
	if err != nil {
		return nil, fmt.Errorf("cannot list Docker containers: %w", err)
	}
	containerNames := make(map[string]struct{})
	for _, name := range utils.StringToLines(dockerPsNames) {
		containerNames[name] = struct{}{}
	}
	for _, name := range requiredContainerNames {
		if _, ok := containerNames[name]; !ok {
			containersNotRunning = append(containersNotRunning, name)
		}
	}
	return containersNotRunning, nil
}

func (d *docker) BadContainers(ctx context.Context) (badStatus []string, err error) {
	dockerPsStatus, err := d.commander.Run(ctx, "docker", "ps", "--format", "'Container {{.Names}} is {{.Status}}'")
	if err != nil {
		return nil, fmt.Errorf("cannot list Docker containers: %w", err)
	}
	allStatus := utils.StringToLines(dockerPsStatus)
	for _, status := range allStatus {
		if strings.Contains(status, "unhealthy") || strings.Contains(status, "restarting") {
			badStatus = append(badStatus, status)
		}
	}
	return badStatus, nil
}
