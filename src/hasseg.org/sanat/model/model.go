package model

import (
	_ "fmt"
)

type TranslationFormatDataType int

const (
	DataTypeNone    TranslationFormatDataType = iota
	DataTypeObject  TranslationFormatDataType = iota
	DataTypeInteger TranslationFormatDataType = iota
	DataTypeString  TranslationFormatDataType = iota
	DataTypeFloat   TranslationFormatDataType = iota
)

type TranslationPlatform int

const (
	PlatformNone    TranslationPlatform = iota
	PlatformApple   TranslationPlatform = iota
	PlatformAndroid TranslationPlatform = iota
	PlatformWindows TranslationPlatform = iota
)

type TranslationValueSegment struct {
	Text               string
	IsFormatSpecifier  bool
	SemanticOrderIndex int
	DataType           TranslationFormatDataType
	NumberOfDecimals   int
}

type TranslationValue struct {
	Language string
	Segments []TranslationValueSegment
}

type Translation struct {
	Key       string
	Values    []TranslationValue
	Platforms []TranslationPlatform
}

type TranslationSection struct {
	Name         string
	Translations []Translation
}

type TranslationSet struct {
	Sections  []TranslationSection
	Languages map[string]bool
}

func NewTranslationSet() TranslationSet {
	return TranslationSet{Languages: make(map[string]bool)}
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

func (translation Translation) IsForPlatform(givenPlatform TranslationPlatform) bool {
	if len(translation.Platforms) == 0 {
		return true
	}
	for _, platform := range translation.Platforms {
		if platform == givenPlatform {
			return true
		}
	}
	return false
}

func NewTextSegment(text string) TranslationValueSegment {
	return TranslationValueSegment{Text: text}
}

func NewFormatSpecifierSegment(dataType TranslationFormatDataType,
	numDecimals int,
	semanticOrderIndex int) TranslationValueSegment {
	return TranslationValueSegment{
		IsFormatSpecifier:  true,
		SemanticOrderIndex: semanticOrderIndex,
		DataType:           dataType,
		NumberOfDecimals:   numDecimals,
	}
}
