package hardware

import (
	"context"
	"strings"
	"time"

	"github.com/c9s/goprocinfo/linux"
	"github.com/capnm/sysinfo"
	"github.com/qdm12/golibs/command"
)

type Hardware interface {
	CPUPercentUsage() (usage int, err error)
	ProcessesCount() (processes int)
	RAMPercentUsage() (usage int)
	Uptime() time.Duration
	PartitionsUsage(ctx context.Context) (partitionsUsage map[string]int, warnings []string, err error)
	IsZpoolInstalled(ctx context.Context) bool
	GetPools(ctx context.Context) (poolNames []string, err error)
	GetPoolHealth(ctx context.Context, poolName string) (health string, err error)
	GetPoolErrors(ctx context.Context, poolName string) (errors string, err error)
	GetPoolCapacity(ctx context.Context, poolName string) (capacity int, err error)
}

type hardware struct {
	commander      Commander
	dockerRootPath string
	readProcStat   func(path string) (*linux.Stat, error)
	getSysInfo     func() *sysinfo.SI
}

type Commander interface {
	Run(cmd command.ExecCmd) (output string, err error)
}

func New(cmd Commander, dockerRootPath string) Hardware {
	return &hardware{
		commander:      cmd,
		dockerRootPath: dockerRootPath,
		readProcStat:   linux.ReadStat,
		getSysInfo:     sysinfo.Get,
	}
}

func stringToLines(s string) (lines []string) {
	lines = strings.Split(s, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		if len(line) > 0 {
			nonEmptyLines++
			if nonEmptyLines == 2 {
				break
			}
		}
	}
	if nonEmptyLines < 2 {
		return []string{s}
	}
	return lines
}
