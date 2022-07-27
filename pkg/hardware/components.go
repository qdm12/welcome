package hardware

import (
	"fmt"
	"math"
	"time"
)

func (hw *Hardware) CPUPercentUsage() (usage int, err error) {
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
	percentage := 100 * totalUse / (totalUse + totalIdle) //nolint:gomnd
	usage = int(math.Round(percentage*10) / 10)           //nolint:gomnd
	return usage, nil
}

func (hw *Hardware) ProcessesCount() (processes int) {
	return int(hw.getSysInfo().Procs)
}

func (hw *Hardware) RAMPercentUsage() (usage int) {
	systemInfo := hw.getSysInfo()
	free := float64(systemInfo.FreeRam)
	total := float64(systemInfo.TotalRam)
	freePercentage := 100 * free / total          //nolint:gomnd
	usage = int(math.Round(100 - freePercentage)) //nolint:gomnd
	return usage
}

func (hw *Hardware) Uptime() time.Duration {
	return hw.getSysInfo().Uptime
}
