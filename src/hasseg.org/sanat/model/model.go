package model

import (
    _ "fmt"
)

type TranslationValue struct {
    Language string
    Text string
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

func (translation *Translation) AddValue(language string, text string) *TranslationValue {
    translation.Values = append(translation.Values, TranslationValue{Language: language, Text: text})
    return &translation.Values[len(translation.Values)-1]
}
