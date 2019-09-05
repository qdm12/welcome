package terminal

import (
	"os/exec"
	"strings"
	"welcome/pkg/utils"
)

func RunCommand(command string, arg ...string) (output string, err error) {
	cmd := exec.Command(command, arg...)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	output = string(stdout)
	output = strings.TrimSuffix(output, "\n")
	lines := utils.StringToLines(output)
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "'")
		lines[i] = strings.TrimSuffix(lines[i], "'")
	}
	output = strings.Join(lines, "\n")
	return output, nil
}
