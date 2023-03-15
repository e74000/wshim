package wshim

import (
	"fmt"
	"strings"
)

func findValidId(Name, Type string) string {
	var (
		prefix string
		count  int
	)

	for {
		prefix = fmt.Sprintf("wshim-%s-%s-%d", Type, kebabCase(Name), count)
		if !ids[prefix] {
			ids[prefix] = true
			return prefix
		}

		count++
	}
}

func kebabCase(input string) string {
	return strings.Replace(input, " ", "-", -1)
}

func checked(b bool) string {
	if b {
		return "checked"
	} else {
		return "unchecked"
	}
}

// There should be a builtin interface for this lol
func clamp[T int | uint | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](v, min, max T) T {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}
