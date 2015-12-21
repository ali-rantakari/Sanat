package android_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/output/android"
	"hasseg.org/sanat/test"
	"hasseg.org/sanat/util"
)

func TestAndroidFormatSpecifierStringForFormatSpecifier(t *testing.T) {
	val := func(dataType model.TranslationFormatDataType,
		numDecimals int,
		semanticOrderIndex int) string {
		return android.FormatSpecifierStringForFormatSpecifier(model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex))
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
		assert.Equal(t, expected, android.SanitizedForString(input), input)
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

func makeTranslationSet(sectionName string, keyName string, language string, value string) model.TranslationSet {
	ts := model.NewTranslationSet()
	ts.AddSection(sectionName).AddTranslation(keyName).AddValue(language, []model.Segment{model.NewTextSegment(value)})
	return ts
}

func TestOverallXMLFileGeneration(t *testing.T) {
	{
		lang := "en"
		ts := makeTranslationSet("Sektion", "Foo", lang, "Some text")
		x := android.GetStringsFileContents(ts, lang)
		assert.True(t, util.XMLIsValid(x), "")
	}
	{
		lang := "en"
		ts := makeTranslationSet("Sektion -- two dashes", "Foo", lang, "Some text")
		x := android.GetStringsFileContents(ts, lang)
		assert.True(t, util.XMLIsValid(x), "-- in XML comment (section name)")
	}
}

func TestComprehensiveInput(t *testing.T) {
	set := test.GetComprehensiveTestInputTranslationSet()
	for language, _ := range set.Languages {
		output := android.GetStringsFileContents(set, language)
		assert.True(t, util.XMLIsValid(output), language)
	}
}
