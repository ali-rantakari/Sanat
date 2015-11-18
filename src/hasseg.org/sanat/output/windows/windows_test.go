package windows_test

import (
	"bytes"
	"encoding/xml"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/output/windows"
)

func TestWindowsFormatSpecifierStringForFormatSpecifier(t *testing.T) {
	val := func(dataType model.TranslationFormatDataType,
		index int,
		numDecimals int,
		semanticOrderIndex int) string {
		return windows.FormatSpecifierStringForFormatSpecifier(model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex), index)
	}

	// Data types
	assert.Equal(t, "{0}", val(model.DataTypeObject, 0, -1, -1), "")
	assert.Equal(t, "{0}", val(model.DataTypeString, 0, -1, -1), "")
	assert.Equal(t, "{0}", val(model.DataTypeFloat, 0, -1, -1), "")
	assert.Equal(t, "{0}", val(model.DataTypeInteger, 0, -1, -1), "")

	// Semantic order index
	assert.Equal(t, "{0}", val(model.DataTypeFloat, 0, -1, 0), "")
	assert.Equal(t, "{77}", val(model.DataTypeFloat, 77, -1, 0), "")
	assert.Equal(t, "{0}", val(model.DataTypeFloat, 0, -1, 1), "Output index is 0-based while explicit order index is 1-based")
	assert.Equal(t, "{11}", val(model.DataTypeFloat, 0, -1, 12), "Output index is 0-based while explicit order index is 1-based")
	assert.Equal(t, "{11}", val(model.DataTypeFloat, 77, -1, 12), "Explicit order index overrides actual index")

	// Decimal count
	assert.Equal(t, "{0:F0}", val(model.DataTypeFloat, 0, 0, -1), "")
	assert.Equal(t, "{0:F1}", val(model.DataTypeFloat, 0, 1, -1), "")
	assert.Equal(t, "{0:F34}", val(model.DataTypeFloat, 0, 34, -1), "")
	assert.Equal(t, "{2:F1}", val(model.DataTypeFloat, 0, 1, 3), "Decimal count together with semantic order index")
	assert.Equal(t, "{0}", val(model.DataTypeObject, 0, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "{0}", val(model.DataTypeString, 0, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "{0}", val(model.DataTypeInteger, 0, 1, -1), "Decimal count is only for floats")
}

func TestTextSanitizedForWindowsString(t *testing.T) {
	ass := func(expected string, input string) {
		assert.Equal(t, expected, windows.SanitizedForStringValue(input), input)
	}

	ass("", "")
	ass("Foo", "Foo")

	// XML-escaping
	ass("&lt;Foo&gt;", "<Foo>")
}

func TestTextSanitizedForWindowsResourceName(t *testing.T) {
	ass := func(expected string, input string) {
		assert.Equal(t, expected, windows.SanitizedForKey(input), input)
	}

	ass("", "")
	ass("Foo", "Foo")

	// XML-escaping
	ass("&lt;Foo&gt;", "<Foo>")

	// Keys must be valid C# identifiers — assert that we handle at least
	// some common cases:
	ass("Foo_Bar", "Foo.Bar")
	ass("Foo_Bar", "Foo Bar")
}

func xmlIsValid(xmlString string) bool {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlString))
	for {
		_, err := decoder.Token()
		if err == nil {
			continue
		} else if err == io.EOF {
			break
		}
		return false
	}
	return true
}

func makeTranslationSet(sectionName string, keyName string, language string, value string) model.TranslationSet {
	ts := model.NewTranslationSet()
	ts.AddSection(sectionName).AddTranslation(keyName).AddValue(language, []model.Segment{model.NewTextSegment(value)})
	return ts
}

func TestOverallXMLFileGeneration(t *testing.T) {
	{
		lang := "en"
		ts := makeTranslationSet("Sektion", "Foo", lang, "Some text")
		x := windows.GetStringsFileContents(ts, lang)
		assert.True(t, xmlIsValid(x), "")
	}
	{
		lang := "en"
		ts := makeTranslationSet("Sektion -- two dashes", "Foo", lang, "Some text")
		x := windows.GetStringsFileContents(ts, lang)
		assert.True(t, xmlIsValid(x), "-- in XML comment (section name)")
	}
}
