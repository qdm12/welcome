package hardware

import (
	"strings"

	"github.com/c9s/goprocinfo/linux"
	"github.com/capnm/sysinfo"
	"github.com/qdm12/golibs/command"
)

type Hardware struct {
	commander      Commander
	dockerRootPath string
	readProcStat   func(path string) (*linux.Stat, error)
	getSysInfo     func() *sysinfo.SI
}

type Commander interface {
	Run(cmd command.ExecCmd) (output string, err error)
}

func New(cmd Commander, dockerRootPath string) *Hardware {
	return &Hardware{
		commander:      cmd,
		dockerRootPath: dockerRootPath,
		readProcStat:   linux.ReadStat,
		getSysInfo:     sysinfo.Get,
	}
}

func stringToLines(s string) (lines []string) {
	lines = strings.Split(s, "\n")

	const minNonEmptyLines = 2
	nonEmptyLines := 0

	for _, line := range lines {
		if len(line) > 0 {
			nonEmptyLines++
			if nonEmptyLines == minNonEmptyLines {
				break
			}
		}
	}
	if nonEmptyLines < minNonEmptyLines {
		return []string{s}
	}
	return lines
}
