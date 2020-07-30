package hardware

import (
	"context"
	"fmt"
	"strconv"
)

func (hw *hardware) IsZpoolInstalled(ctx context.Context) bool {
	_, err := hw.cmd.Run(ctx, "zpool", "version")
	return err == nil
}

func (hw *hardware) GetPools(ctx context.Context) (poolNames []string, err error) {
	output, err := hw.cmd.Run(ctx, "zpool", "list", "-o", "name", "-H")
	if err != nil {
		return nil, fmt.Errorf("cannot list zpools: %w", err)
	}
	poolNames = stringToLines(output)
	return poolNames, nil
}

func (hw *hardware) GetPoolHealth(ctx context.Context, poolName string) (health string, err error) {
	health, err = hw.cmd.Run(ctx, "zpool", "list", poolName, "-o", "health", "-H")
	if err != nil {
		return "", fmt.Errorf("cannot get zpool %s health: %w", poolName, err)
	}
	if health == "ONLINE" {
		return "", nil
	}
	return health, nil
}

func (hw *hardware) GetPoolErrors(ctx context.Context, poolName string) (errors string, err error) {
	poolStatus, err := hw.cmd.Run(ctx, "zpool", "status", poolName)
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

func (hw *hardware) GetPoolCapacity(ctx context.Context, poolName string) (capacity int, err error) {
	output, err := hw.cmd.Run(ctx, "zpool", "list", poolName, "-o", "capacity", "-H")
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
