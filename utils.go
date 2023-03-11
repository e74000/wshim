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
		if ids[prefix] {
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
