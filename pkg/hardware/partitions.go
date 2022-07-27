package hardware

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func (hw *Hardware) PartitionsUsage(ctx context.Context) (
	partitionsUsage map[string]int, warnings []string, err error) {
	lines, err := hw.getDrivesRawMetadata(ctx)
	if err != nil {
		return nil, nil, err
	}
	partitions, warnings := hw.processPartitionsRawMetadata(lines)
	partitionsUsage = make(map[string]int, len(partitions))
	for _, partition := range partitions {
		partitionsUsage[partition.filesystem] = partition.use
	}
	return partitionsUsage, warnings, nil
}

func (hw *Hardware) getDrivesRawMetadata(ctx context.Context) (lines []string, err error) {
	cmd := exec.CommandContext(ctx, "df", "-T")
	output, err := hw.commander.Run(cmd)
	if err != nil {
		return nil, fmt.Errorf("cannot get drives raw metadata: %w", err)
	}
	lines = stringToLines(output)
	return lines[1:], nil
}

type partitionData struct {
	filesystem    string
	partitionType string
	use           int
	mountedOn     string
}

func (hw *Hardware) processPartitionsRawMetadata(lines []string) (partitions []partitionData, warnings []string) {
	for _, line := range lines {
		if partitionRawDataShouldBeSkipped(line, []string{hw.dockerRootPath}) {
			continue
		}
		partition, err := makePartitionData(line)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("Cannot extract partition data from %q: %s", line, err))
			continue
		}
		partitions = append(partitions, partition)
	}
	return partitions, warnings
}

var (
	ErrExtractPartitionInformation = errors.New("cannot extract partition information")
)

func makePartitionData(line string) (data partitionData, err error) {
	columns := strings.Fields(line)
	const minColumns = 7
	if len(columns) < minColumns {
		return data, fmt.Errorf("%w: %q has less than %d columns",
			ErrExtractPartitionInformation, line, minColumns)
	}
	data.filesystem = columns[0]
	data.partitionType = columns[1]
	percent := strings.TrimSuffix(columns[5], "%")
	data.use, err = strconv.Atoi(percent)
	if err != nil {
		return data, fmt.Errorf("cannot extract partition usage percent: %w", err)
	}
	data.mountedOn = columns[6]
	return data, nil
}

func partitionRawDataShouldBeSkipped(line string, ignoredMountPoints []string) (skip bool) {
	CIFSEncryptedShare := strings.HasPrefix(line, "//")
	isBootMountpoint := strings.HasSuffix(line, "/boot/efi") ||
		strings.HasSuffix(line, "/boot") ||
		strings.HasSuffix(line, "/efi")
	isCIFS := strings.Contains(line, " cifs ")
	switch {
	case len(line) == 0, CIFSEncryptedShare, isBootMountpoint, isCIFS:
		return true
	}
	columns := strings.Fields(line)
	partitionType := ""
	if len(columns) > 1 {
		partitionType = columns[1]
	}
	ignoredPartitionTypes := []string{"zfs", "devtmpfs", "tmpfs", "cifs", "overlay"}
	for _, t := range ignoredPartitionTypes {
		if partitionType == t {
			return true
		}
	}
	mountpoint := columns[len(columns)-1]
	for i := range ignoredMountPoints {
		if strings.Contains(mountpoint, ignoredMountPoints[i]) {
			return true
		}
	}
	return false
}
