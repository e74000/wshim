package wshim

import (
	"fmt"
	"strings"
)

// findValidId returns a valid id for a given name and type.
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

// kebabCase converts a string to kebab-case.
func kebabCase(input string) string {
	return strings.Replace(input, " ", "-", -1)
}

// checked formats a boolean for checkbox elements
func checked(b bool) string {
	if b {
		return "checked"
	} else {
		return "unchecked"
	}
}

// clamp clamps a value between min and max.
func clamp[T int | uint | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64](v, min, max T) T {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}
