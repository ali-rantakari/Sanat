package util_test

import (
	"github.com/stretchr/testify/assert"
	"hasseg.org/sanat/util"
	"testing"
)

func TestComponentsFromCommaSeparatedList(t *testing.T) {
	ass := func(expected []string, given string) {
		assert.Equal(t, expected, util.ComponentsFromCommaSeparatedList(given), given)
	}

	ass([]string{}, "")
	ass([]string{"", ""}, ",")
	ass([]string{"a", "b"}, "a,b")

	// Trims whitespace
	ass([]string{"a", "b"}, " a , b ")
}
