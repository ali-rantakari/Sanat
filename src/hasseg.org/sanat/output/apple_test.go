package output_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "hasseg.org/sanat/output"
    "hasseg.org/sanat/model"
)

func TestAppleFormatSpecifierStringForFormatSpecifier(t *testing.T) {
    val := func(dataType model.TranslationFormatDataType,
                numDecimals int,
                semanticOrderIndex int) string {
        return output.AppleFormatSpecifierStringForFormatSpecifier(model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex))
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
        assert.Equal(t, expected, output.TextSanitizedForAppleString(input), input)
    }

    ass("", "")
    ass("Foo", "Foo")
    ass("Per %% cent", "Per % cent")
    ass("Foo %%@", "Foo %@")
    ass("Foo %%@ %%@ %%@", "Foo %@ %@ %@")
}
