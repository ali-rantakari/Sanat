package model

import ()

// TranslationFormatDataType is the “enum” type for format
// specifier data types.
type TranslationFormatDataType int

const (
	DataTypeNone TranslationFormatDataType = iota
	DataTypeObject
	DataTypeInteger
	DataTypeString
	DataTypeFloat
)

// TranslationPlatform is the “enum” type for software
// platform indicator values.
type TranslationPlatform int

const (
	PlatformNone TranslationPlatform = iota
	PlatformApple
	PlatformAndroid
	PlatformWindows
)

// TextSegment is a piece of a translation string value
// containing text.
type TextSegment struct {
	Text string
}

// FormatSpecifierSegment is a piece of a translation string
// value that describes the attributes of a format specifier.
type FormatSpecifierSegment struct {
	SemanticOrderIndex int
	DataType           TranslationFormatDataType
	NumberOfDecimals   int
}

// Segment is the “union” type for translation string value
// segments.
type Segment interface{}

// TranslationValue is a value for a specific language for
// a translation string.
type TranslationValue struct {
	Language string
	Segments []Segment
}

// Translation is a unique localizable string containing
// values for N languages. It can be limited only to specific
// platforms.
type Translation struct {
	Key       string
	Values    []TranslationValue
	Platforms []TranslationPlatform
	Tags      []string
	Comment   string
}

// TranslationSection is a named group of Translations.
type TranslationSection struct {
	Name         string
	Translations []Translation
}

// TranslationSet is a set of TranslationSections.
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

func (translation *Translation) AddValue(language string, segments []Segment) *TranslationValue {
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

func NewTextSegment(text string) TextSegment {
	return TextSegment{Text: text}
}

func NewFormatSpecifierSegment(dataType TranslationFormatDataType,
	numDecimals int,
	semanticOrderIndex int) FormatSpecifierSegment {
	return FormatSpecifierSegment{
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
