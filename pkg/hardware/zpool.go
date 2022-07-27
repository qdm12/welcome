package hardware

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
)

func (hw *Hardware) IsZpoolInstalled(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "zpool", "version")
	_, err := hw.commander.Run(cmd)
	return err == nil
}

func (hw *Hardware) GetPools(ctx context.Context) (poolNames []string, err error) {
	cmd := exec.CommandContext(ctx, "zpool", "list", "-o", "name", "-H")
	output, err := hw.commander.Run(cmd)
	if err != nil {
		return nil, fmt.Errorf("cannot list zpools: %w", err)
	}
	poolNames = stringToLines(output)
	return poolNames, nil
}

func (hw *Hardware) GetPoolHealth(ctx context.Context, poolName string) (health string, err error) {
	cmd := exec.CommandContext(ctx, "zpool", "list", poolName, "-o", "health", "-H")
	health, err = hw.commander.Run(cmd)
	if err != nil {
		return "", fmt.Errorf("cannot get zpool %s health: %w", poolName, err)
	}
	if health == "ONLINE" {
		return "", nil
	}
	return health, nil
}

func (hw *Hardware) GetPoolErrors(ctx context.Context, poolName string) (errors string, err error) {
	cmd := exec.CommandContext(ctx, "zpool", "status", poolName)
	poolStatus, err := hw.commander.Run(cmd)
	if err != nil {
		return "", fmt.Errorf("cannot get zpool %s status: %w", poolName, err)
	}
	lines := stringToLines(poolStatus)
	errorsLine := lines[len(lines)-1]
	errors = errorsLine[8:]
	if errors == "No known data errors" {
		return "", nil
	}
	return errors, nil
}

func (hw *Hardware) GetPoolCapacity(ctx context.Context, poolName string) (capacity int, err error) {
	cmd := exec.CommandContext(ctx, "zpool", "list", poolName, "-o", "capacity", "-H")
	output, err := hw.commander.Run(cmd)
	if err != nil {
		return 0, fmt.Errorf("cannot get zpool %s capacity: %w", poolName, err)
	}
	output = output[:len(output)-1] // removes % sign
	capacity, err = strconv.Atoi(output)
	if err != nil {
		return 0, fmt.Errorf("cannot get zpool %s capacity: %w", poolName, err)
	}
	return capacity, nil
}
