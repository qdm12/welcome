package hardware

import (
	"fmt"
	"strconv"
	"strings"
	"welcome/pkg/display"
	"welcome/pkg/terminal"
	"welcome/pkg/utils"
)

func GetPartitionsUsage() (partitionsUsage map[string]int, err error) {
	lines, err := getDrivesRawMetadata()
	if err != nil {
		return nil, err
	}
	partitionsUsage = make(map[string]int)
	data := processPartitionsRawMetadata(lines)
	for i := range data {
		partitionsUsage[data[i].filesystem] = data[i].use
	}
	return partitionsUsage, nil
}

func getDrivesRawMetadata() (lines []string, err error) {
	output, err := terminal.RunCommand("df", "-T")
	if err != nil {
		return nil, err
	}
	lines = utils.StringToLines(output)
	return lines[1:len(lines)], nil
}

type partitionData struct {
	filesystem    string
	partitionType string
	use           int
	mountedOn     string
}

func processPartitionsRawMetadata(lines []string) (data []partitionData) {
	for _, line := range lines {
		if partitionRawDataShouldBeSkipped(line) {
			continue
		}
		d, err := makePartitionData(line)
		if err != nil {
			display.Error("%s", err)
			continue
		}
		data = append(data, d)
	}
	return data
}

func makePartitionData(line string) (data partitionData, err error) {
	columns := strings.Fields(line)
	if len(columns) < 7 {
		return data, fmt.Errorf("%s has less than 7 columns", line)
	}
	data.filesystem = columns[0]
	data.partitionType = columns[1]
	percent := strings.TrimSuffix(columns[5], "%")
	data.use, err = strconv.Atoi(percent)
	if err != nil {
		return data, err
	}
	data.mountedOn = columns[6]
	return data, nil
}

func partitionRawDataShouldBeSkipped(line string) (skip bool) {
	if len(line) == 0 {
		return true
	} else if strings.HasPrefix(line, "//") { // CIFS encryted share
		return true
	} else if strings.HasSuffix(line, "/boot/efi") { // ignore EFI mountpoint
		return true
	}
	columns := strings.Fields(line)
	partitionType := ""
	if len(columns) > 1 {
		partitionType = columns[1]
	}
	ignoredPartitionTypes := []string{"zfs", "devtmpfs", "tmpfs", "cifs", "overlay", "zfs"}
	for _, t := range ignoredPartitionTypes {
		if partitionType == t {
			return true
		}
	}
	mountpoint := columns[len(columns)-1]
	if strings.Contains(mountpoint, "/var/lib/docker") {
		return true
	}
	return false
}
