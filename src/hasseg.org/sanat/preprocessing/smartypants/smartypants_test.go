package smartypants_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hasseg.org/sanat/preprocessing/smartypants"
)

func TestSmartypantsProcessing(t *testing.T) {
	pp := smartypants.Preprocessor{}

	ass := func(expected string, input string) {
		assert.Equal(t, expected, pp.ProcessRawValue(input), input)
	}

	ass("", "")
	ass("Hello", "Hello")

	// Basic Smartypants substitutions on individual word
	ass("‘Hello’", "'Hello'")
	ass("“Hello”", "\"Hello\"")
	ass("Ali’s", "Ali's")
	ass("Alis’", "Alis'")

	// Corner cases (enable when we're ready to tackle these)
	//ass("My 50\" TV", "My 50\" TV")
	//ass("50\" TVs", "50\" TVs")

	// Preserves leading/trailing whitespace
	ass("\t‘Hello’\t", "\t'Hello'\t")
	ass("  ‘Hello’  ", "  'Hello'  ")
}
