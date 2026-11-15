package resp

import (
	"strconv"
	"strings"
)

// ParseRequest turns a raw request like "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
// into its parts: ["ECHO", "hey"].
func ParseRequest(data string) []string {
	lines := strings.Split(data, "\r\n")
	if len(lines[0]) == 0 {
		return nil
	}

	switch lines[0][0] {
	case '$':
		if len(lines) < 2 {
			return nil
		}
		return []string{lines[1]}
	case '*':
		return parseArray(lines)
	}
	return nil
}

func parseArray(lines []string) []string {
	count, err := strconv.Atoi(lines[0][1:])
	if err != nil {
		return nil
	}

	parts := []string{}
	for i := 0; i < count; i++ {
		// each element takes two lines: "$<length>" and then the value
		valueIndex := 2 + i*2
		if valueIndex >= len(lines) {
			break
		}
		parts = append(parts, lines[valueIndex])
	}
	return parts
}
