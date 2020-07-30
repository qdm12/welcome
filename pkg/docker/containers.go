package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
)

func (d *docker) CountContainers(ctx context.Context) (count int, err error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return 0, fmt.Errorf("cannot list Docker containers: %w", err)
	}
	return len(containers), err
}

func (d *docker) AreContainerRunning(ctx context.Context, requiredContainerNames []string) (containersNotRunning []string, err error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot list Docker containers: %w", err)
	}
	containerNames := make(map[string]struct{})
	for _, container := range containers {
		for _, name := range container.Names {
			containerNames[strings.TrimPrefix(name, "/")] = struct{}{}
		}
	}
	for _, name := range requiredContainerNames {
		if _, ok := containerNames[name]; !ok {
			containersNotRunning = append(containersNotRunning, name)
		}
	}
	return containersNotRunning, nil
}

func (d *docker) BadContainers(ctx context.Context) (containerNameToState map[string]string, err error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot list Docker containers: %w", err)
	}
	containerNameToState = map[string]string{}
	for _, container := range containers {
		containerName := strings.TrimPrefix(container.Names[0], "/")
		lowercaseStatus := strings.ToLower(container.Status)
		switch {
		case strings.Contains(lowercaseStatus, "unhealthy"):
			containerNameToState[containerName] = "unhealthy"
		case strings.Contains(lowercaseStatus, "restarting"):
			containerNameToState[containerName] = "restarting"
		case !strings.HasPrefix(lowercaseStatus, "up "):
			containerNameToState[containerName] = lowercaseStatus
		}
	}
	return containerNameToState, nil
}
