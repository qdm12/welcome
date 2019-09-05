package hardware

import (
	"strconv"
	"welcome/pkg/terminal"
	"welcome/pkg/utils"
)

func IsZpoolInstalled() bool {
	_, err := terminal.RunCommand("zpool", "version")
	if err != nil {
		return false
	}
	return true
}

func GetPools() (poolNames []string, err error) {
	output, err := terminal.RunCommand("zpool", "list", "-o", "name", "-H")
	if err != nil {
		return nil, err
	}
	poolNames = utils.StringToLines(output)
	return poolNames, nil
}

func GetPoolHealth(poolName string) (health string, err error) {
	health, err = terminal.RunCommand("zpool", "list", poolName, "-o", "health", "-H")
	if err != nil {
		return "", err
	}
	if health == "ONLINE" {
		return "", nil
	}
	return health, nil
}

func GetPoolErrors(poolName string) (errors string, err error) {
	poolStatus, err := terminal.RunCommand("zpool", "status", poolName)
	if err != nil {
		return "", err
	}
	lines := utils.StringToLines(poolStatus)
	errorsLine := lines[len(lines)-1]
	errors = errorsLine[8:]
	if errors == "No known data errors" {
		return "", nil
	}
	return errors, nil
}

func GetPoolCapacity(poolName string) (capacity int, err error) {
	output, err := terminal.RunCommand("zpool", "list", poolName, "-o", "capacity", "-H")
	if err != nil {
		return 0, err
	}
	output = output[:len(output)-1] // removes % sign
	capacity, err = strconv.Atoi(output)
	return capacity, err
}
