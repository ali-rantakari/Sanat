package java_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/output/java"
	"hasseg.org/sanat/test"
	"hasseg.org/sanat/util"
)

func TestJavaFormatSpecifierStringForFormatSpecifier(t *testing.T) {
	val := func(dataType model.TranslationFormatDataType,
		index int,
		numDecimals int,
		semanticOrderIndex int) string {
		return java.FormatSpecifierStringForFormatSpecifier(model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex), index)
	}

	// Data types
	assert.Equal(t, "{0}", val(model.DataTypeObject, 0, -1, -1), "")
	assert.Equal(t, "{0}", val(model.DataTypeString, 0, -1, -1), "")
	assert.Equal(t, "{0,number}", val(model.DataTypeFloat, 0, -1, -1), "")
	assert.Equal(t, "{0,number,integer}", val(model.DataTypeInteger, 0, -1, -1), "")

	// Semantic order index
	assert.Equal(t, "{0}", val(model.DataTypeObject, 0, -1, 0), "")
	assert.Equal(t, "{77}", val(model.DataTypeObject, 77, -1, 0), "")
	assert.Equal(t, "{0}", val(model.DataTypeObject, 0, -1, 1), "Output index is 0-based while explicit order index is 1-based")
	assert.Equal(t, "{11}", val(model.DataTypeObject, 0, -1, 12), "Output index is 0-based while explicit order index is 1-based")
	assert.Equal(t, "{11}", val(model.DataTypeObject, 77, -1, 12), "Explicit order index overrides actual index")

	// Decimal count
	assert.Equal(t, "{0,number,#}", val(model.DataTypeFloat, 0, 0, -1), "")
	assert.Equal(t, "{0,number,#.#}", val(model.DataTypeFloat, 0, 1, -1), "")
	assert.Equal(t, "{0,number,#.#####}", val(model.DataTypeFloat, 0, 5, -1), "")
	assert.Equal(t, "{2,number,#.#}", val(model.DataTypeFloat, 0, 1, 3), "Decimal count together with semantic order index")
	assert.Equal(t, "{0}", val(model.DataTypeObject, 0, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "{0}", val(model.DataTypeString, 0, 1, -1), "Decimal count is only for floats")
	assert.Equal(t, "{0,number,integer}", val(model.DataTypeInteger, 0, 1, -1), "Decimal count is only for floats")
}

func TestTextSanitizedForJavaString(t *testing.T) {
	ass := func(expected string, input string) {
		assert.Equal(t, expected, java.SanitizedForStringValue(input), input)
	}

	ass("", "")
	ass("Foo", "Foo")

	// https://docs.oracle.com/javase/7/docs/api/java/text/MessageFormat.html

	// Escaping curly braces with single quotes
	//
	ass(util.XMLEscaped("eka '{}' toka"), "eka {} toka")
	ass(util.XMLEscaped("eka '{{}}' toka"), "eka {{}} toka")
	ass(util.XMLEscaped("eka '{' toka"), "eka { toka")
	ass(util.XMLEscaped("eka '{'0'}' toka"), "eka {0} toka")
	ass(util.XMLEscaped("eka '{' keski moro '}' toka"), "eka { keski moro } toka")

	// Escaping single quotes themselves
	//
	ass(util.XMLEscaped("eka '' toka"), "eka ' toka")
	ass(util.XMLEscaped("eka '''{}''' toka"), "eka '{}' toka")
	ass(util.XMLEscaped("eka '{''}' toka"), "eka {'} toka")
	ass(util.XMLEscaped("eka '{''''}' toka"), "eka {''} toka")

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
		x := java.GetPropertiesFileContents(ts, lang)
		assert.True(t, util.XMLIsValid(x), "")
	}
	{
		lang := "en"
		ts := makeTranslationSet("Sektion -- two dashes", "Foo", lang, "Some text")
		x := java.GetPropertiesFileContents(ts, lang)
		assert.True(t, util.XMLIsValid(x), "-- in XML comment (section name)")
	}
}

func TestComprehensiveInput(t *testing.T) {
	set := test.GetComprehensiveTestInputTranslationSet()
	for language, _ := range set.Languages {
		output := java.GetPropertiesFileContents(set, language)
		assert.True(t, util.XMLIsValid(output), language)
	}
}
