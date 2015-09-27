package util

import (
	"strings"
)

func ComponentsFromCommaSeparatedList(text string) []string {
	ret := make([]string, 0)
	for _, s := range strings.Split(text, ",") {
		ret = append(ret, strings.TrimSpace(s))
	}
	return ret
}
