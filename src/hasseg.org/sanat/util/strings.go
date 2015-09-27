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

func LeadingWhitespace(s string) string {
	ret := ""
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			ret += " "
		} else {
			break
		}
	}
	return ret
}

func TrailingWhitespace(s string) string {
	ret := ""
	for i := len(s) - 1; 0 <= i; i-- {
		if s[i] == ' ' {
			ret += " "
		} else {
			break
		}
	}
	return ret
}