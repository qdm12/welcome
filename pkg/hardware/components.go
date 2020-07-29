package hardware

import (
	"fmt"
	"math"
	"time"
)

func (hw *hardware) CPUPercentUsage() (usage int, err error) {
	stat, err := hw.readProcStat("/proc/stat")
	if err != nil {
		return 0, fmt.Errorf("cannot get CPU percent usage: %w", err)
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

func (hw *hardware) ProcessesCount() (processes int) {
	return int(hw.getSysInfo().Procs)
}

func (hw *hardware) RAMPercentUsage() (usage int) {
	systemInfo := hw.getSysInfo()
	free := float64(systemInfo.FreeRam)
	total := float64(systemInfo.TotalRam)
	freePercentage := 100 * free / total
	usage = int(math.Round(100 - freePercentage))
	return usage
}

func (hw *hardware) Uptime() time.Duration {
	return hw.getSysInfo().Uptime
}
