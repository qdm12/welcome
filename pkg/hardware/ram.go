package hardware

import (
	"math"

	"github.com/capnm/sysinfo"
)

func GetRAMUsage() (usage int) {
	systemInfo := sysinfo.Get()
	free := float64(systemInfo.FreeRam)
	total := float64(systemInfo.TotalRam)
	freePercentage := 100 * free / total
	usage = int(math.Round(100 - freePercentage))
	return usage
}
