package util

import (
	"strings"
)

func ComponentsFromCommaSeparatedList(text string) []string {
	ret := make([]string, 0)
	if len(strings.TrimSpace(text)) == 0 {
		return ret
	}
	for _, s := range strings.Split(text, ",") {
		ret = append(ret, strings.TrimSpace(s))
	}
	return ret
}
