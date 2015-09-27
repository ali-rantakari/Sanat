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

func TestLeadingWhitespace(t *testing.T) {
	ass := func(expected string, given string) {
		assert.Equal(t, expected, util.LeadingWhitespace(given), given)
	}

	ass("", "")
	ass("", "Moro")
	ass(" ", " Moro")
	ass("   ", "   Moro")
	ass("", ".   Moro")
	ass("", "xx   Moro")
	ass("", "Moro ")
	ass("", "Moro  ")
	ass("", "Moro  xx")
}

func TestTrailingWhitespace(t *testing.T) {
	ass := func(expected string, given string) {
		assert.Equal(t, expected, util.TrailingWhitespace(given), given)
	}

	ass("", "")
	ass("", "Moro")
	ass("", " Moro")
	ass("", "   Moro")
	ass("", ".   Moro")
	ass("", "xx   Moro")
	ass(" ", "Moro ")
	ass("  ", "Moro  ")
	ass("", "Moro  xx")
}
