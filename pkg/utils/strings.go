package utils

import "strings"

func StringToLines(s string) (lines []string) {
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
