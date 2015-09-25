package parser_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "hasseg.org/sanat/parser"
    "hasseg.org/sanat/model"
)

func TestNewFormatSpecifierSegmentFromSpecifierText(t *testing.T) {
    val := parser.NewFormatSpecifierSegmentFromSpecifierText
    seg := model.NewFormatSpecifierSegment

    ass := func(f string, s model.TranslationValueSegment) {
        assert.Equal(t, s, val(f), f)
    }

    // Data types
    ass("{@}", seg(model.DataTypeObject, -1, -1))
    ass("{d}", seg(model.DataTypeInteger, -1, -1))
    ass("{f}", seg(model.DataTypeFloat, -1, -1))
    ass("{s}", seg(model.DataTypeString, -1, -1))

    // Semantic order index
    assert.Equal(t, seg(model.DataTypeObject, -1, -1), val("{0:@}"), "Order index 0 = none at all")
    assert.Equal(t, seg(model.DataTypeObject, -1, 1), val("{1:@}"), "")
    assert.Equal(t, seg(model.DataTypeObject, -1, 72), val("{72:@}"), "")

    // Decimal count
    ass("{f.}", seg(model.DataTypeFloat, -1, -1))
    ass("{f.0}", seg(model.DataTypeFloat, 0, -1))
    ass("{f.2}", seg(model.DataTypeFloat, 2, -1))

    // Combination
    ass("{1:f.2}", seg(model.DataTypeFloat, 2, 1))
}
