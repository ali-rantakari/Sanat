package output_test

import (
	"bytes"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"hasseg.org/sanat/model"
	"hasseg.org/sanat/output"
	"io"
	"testing"
)

func TestAndroidFormatSpecifierStringForFormatSpecifier(t *testing.T) {
	val := func(dataType model.TranslationFormatDataType,
		numDecimals int,
		semanticOrderIndex int) string {
		return output.AndroidFormatSpecifierStringForFormatSpecifier(model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex))
	}

	// Data types
	assert.Equal(t, "%s", val(model.DataTypeObject, -1, -1), "")
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
	assert.Equal(t, "%s", val(model.DataTypeObject, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "%s", val(model.DataTypeString, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "%d", val(model.DataTypeInteger, 1, -1), "Decimal count is only for floats")
}

func TestTextSanitizedForAndroidString(t *testing.T) {
	ass := func(expected string, input string) {
		assert.Equal(t, expected, output.TextSanitizedForAndroidString(input), input)
	}

	ass("", "")
	ass("Foo", "Foo")

	// Escaping %
	ass("Per %% cent", "Per % cent")
	ass("Foo %%@", "Foo %@")
	ass("Foo %%@ %%@ %%@", "Foo %@ %@ %@")

	// XML-escaping
	ass("&lt;Foo&gt;", "<Foo>")
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
	ts.AddSection(sectionName).AddTranslation(keyName).AddValue(language, []model.TranslationValueSegment{model.NewTextSegment(value)})
	return ts
}

func TestOverallXMLFileGeneration(t *testing.T) {
	{
		lang := "en"
		ts := makeTranslationSet("Sektion", "Foo", lang, "Some text")
		x := output.GetAndroidStringsFileContents(ts, lang)
		assert.True(t, xmlIsValid(x), "")
	}
}
