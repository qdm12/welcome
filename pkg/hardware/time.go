package hardware

import (
	"fmt"

	"github.com/capnm/sysinfo"
)

func GetUptime() string {
	systemInfo := sysinfo.Get()
	return fmt.Sprintf("%s", systemInfo.Uptime)
}
