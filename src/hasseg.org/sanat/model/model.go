package model

import ()

type TranslationFormatDataType int

const (
	DataTypeNone TranslationFormatDataType = iota
	DataTypeObject
	DataTypeInteger
	DataTypeString
	DataTypeFloat
)

type TranslationPlatform int

const (
	PlatformNone TranslationPlatform = iota
	PlatformApple
	PlatformAndroid
	PlatformWindows
)

type TranslationValueTextSegment struct {
	Text string
}

type TranslationValueFormatSpecifierSegment struct {
	SemanticOrderIndex int
	DataType           TranslationFormatDataType
	NumberOfDecimals   int
}

type TranslationValueSegment interface{}

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

func NewTextSegment(text string) TranslationValueTextSegment {
	return TranslationValueTextSegment{Text: text}
}

func NewFormatSpecifierSegment(dataType TranslationFormatDataType,
	numDecimals int,
	semanticOrderIndex int) TranslationValueFormatSpecifierSegment {
	return TranslationValueFormatSpecifierSegment{
		SemanticOrderIndex: semanticOrderIndex,
		DataType:           dataType,
		NumberOfDecimals:   numDecimals,
	}
}

type TranslationValueHandler func(*TranslationValue)

func (set *TranslationSet) IterateTranslationValues(handler TranslationValueHandler) {
	for s := 0; s < len(set.Sections); s++ {
		for t := 0; t < len(set.Sections[s].Translations); t++ {
			for v := 0; v < len(set.Sections[s].Translations[t].Values); v++ {
				handler(&set.Sections[s].Translations[t].Values[v])
			}
		}
	}
}
