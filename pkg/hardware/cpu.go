package hardware

import (
	"math"

	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/capnm/sysinfo"
)

func GetCPUUsage() (usage int, err error) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return 0, err
	}
	var totalUse float64
	var totalIdle float64
	for _, CPUStat := range stat.CPUStats {
		totalUse += float64(CPUStat.User + CPUStat.System)
		totalIdle += float64(CPUStat.Idle)
	}
	percentage := 100 * totalUse / (totalUse + totalIdle)
	usage = int(math.Round(percentage*10) / 10)
	return usage, nil
}

func GetProcesses() (processes int) {
	systemInfo := sysinfo.Get()
	return int(systemInfo.Procs)
}
