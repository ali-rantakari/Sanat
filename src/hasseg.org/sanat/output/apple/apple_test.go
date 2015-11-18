package apple_test

import (
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/output/apple"
	"hasseg.org/sanat/test"
)

func TestAppleFormatSpecifierStringForFormatSpecifier(t *testing.T) {
	val := func(dataType model.TranslationFormatDataType,
		numDecimals int,
		semanticOrderIndex int) string {
		return apple.FormatSpecifierStringForFormatSpecifier(model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex))
	}

	// Data types
	assert.Equal(t, "%@", val(model.DataTypeObject, -1, -1), "")
	assert.Equal(t, "%s", val(model.DataTypeString, -1, -1), "")
	assert.Equal(t, "%f", val(model.DataTypeFloat, -1, -1), "")
	assert.Equal(t, "%d", val(model.DataTypeInteger, -1, -1), "")

	// Semantic order index
	assert.Equal(t, "%f", val(model.DataTypeFloat, -1, 0), "")
	assert.Equal(t, "%1$f", val(model.DataTypeFloat, -1, 1), "")
	assert.Equal(t, "%12$f", val(model.DataTypeFloat, -1, 12), "")

	// Decimal count
	assert.Equal(t, "%.0f", val(model.DataTypeFloat, 0, -1), "")
	assert.Equal(t, "%.1f", val(model.DataTypeFloat, 1, -1), "")
	assert.Equal(t, "%.34f", val(model.DataTypeFloat, 34, -1), "")
	assert.Equal(t, "%3$.1f", val(model.DataTypeFloat, 1, 3), "Decimal count together with semantic order index")
	assert.Equal(t, "%@", val(model.DataTypeObject, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "%s", val(model.DataTypeString, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "%d", val(model.DataTypeInteger, 1, -1), "Decimal count is only for floats")
}

func TestTextSanitizedForAppleString(t *testing.T) {
	ass := func(expected string, input string) {
		assert.Equal(t, expected, apple.SanitizedForStringValue(input), input)
	}

	ass("", "")
	ass("Foo", "Foo")

	ass("Per %% cent", "Per % cent")
	ass("Foo %%@", "Foo %@")
	ass("Foo %%@ %%@ %%@", "Foo %@ %@ %@")

	ass("Foo\\\"bar", "Foo\"bar")
}

func isValidPlist(plistString string) bool {
	cmd := exec.Command("/usr/bin/plutil", "-lint", "-")

	cmdStdin, _ := cmd.StdinPipe()

	cmd.Start()
	io.WriteString(cmdStdin, plistString)
	cmdStdin.Close()

	err := cmd.Wait()
	return (err == nil)
}

func TestComprehensiveInput(t *testing.T) {
	set := test.GetComprehensiveTestInputTranslationSet()
	for language, _ := range set.Languages {
		output := apple.GetStringsFileContents(set, language)
		assert.True(t, isValidPlist(output), language)
	}
}
