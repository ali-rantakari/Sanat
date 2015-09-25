package model

import (
    _ "fmt"
)

type TranslationFormatDataType int

const (
    DataTypeObject TranslationFormatDataType = iota
    DataTypeInteger TranslationFormatDataType = iota
    DataTypeString TranslationFormatDataType = iota
    DataTypeFloat TranslationFormatDataType = iota
)

type TranslationValueSegment struct {
    Text string
    IsFormatSpecifier bool
    SemanticOrderIndex int
    DataType TranslationFormatDataType
    NumberOfDecimals int
}

type TranslationValue struct {
    Language string
    Segments []TranslationValueSegment
}

type Translation struct {
    Key string
    Values []TranslationValue
}

type TranslationSection struct {
    Name string
    Translations []Translation
}

type TranslationSet struct {
    Sections []TranslationSection
}

func NewTranslationSet() TranslationSet {
    return TranslationSet{}
}

func (set *TranslationSet) AddSection(name string) *TranslationSection {
    set.Sections = append(set.Sections, TranslationSection{Name: name})
    return &set.Sections[len(set.Sections)-1]
}

func (section *TranslationSection) AddTranslation(key string) *Translation {
    section.Translations = append(section.Translations, Translation{Key: key})
    return &section.Translations[len(section.Translations)-1]
}

func (translation *Translation) AddValue(language string, segments []TranslationValueSegment) *TranslationValue {
    translation.Values = append(translation.Values, TranslationValue{Language: language, Segments: segments})
    return &translation.Values[len(translation.Values)-1]
}

func NewTextSegment(text string) TranslationValueSegment {
    return TranslationValueSegment{Text: text}
}

func NewFormatSpecifierSegment(dataType TranslationFormatDataType,
                               numDecimals int,
                               semanticOrderIndex int) TranslationValueSegment {
    return TranslationValueSegment{
        IsFormatSpecifier: true,
        SemanticOrderIndex: semanticOrderIndex,
        DataType: dataType,
        NumberOfDecimals: numDecimals,
    }
}
