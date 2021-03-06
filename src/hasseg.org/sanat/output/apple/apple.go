package apple

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"hasseg.org/sanat/model"
)

func FormatSpecifierStringForFormatSpecifier(segment model.FormatSpecifierSegment) string {
	ret := "%"
	if 0 < segment.SemanticOrderIndex {
		ret += strconv.Itoa(segment.SemanticOrderIndex) + "$"
	}
	if segment.DataType == model.DataTypeFloat && 0 <= segment.NumberOfDecimals {
		ret += "." + strconv.Itoa(segment.NumberOfDecimals)
	}
	switch segment.DataType {
	case model.DataTypeInteger:
		ret += "d"
	case model.DataTypeFloat:
		ret += "f"
	case model.DataTypeString:
		fallthrough
	case model.DataTypeObject:
		ret += "@"
	}
	return ret
}

func escapedForComment(s string) string {
	return strings.Replace(s, "*/", "* /", -1)
}

func SanitizedForStringValue(text string) string {
	ret := text
	ret = strings.Replace(ret, "%", "%%", -1)
	ret = strings.Replace(ret, "\"", "\\\"", -1)
	return ret
}

func StringFromSegments(segments []model.Segment) string {
	ret := ""
	for _, segment := range segments {
		switch segment.(type) {
		case model.TextSegment:
			ret += SanitizedForStringValue(segment.(model.TextSegment).Text)
		case model.FormatSpecifierSegment:
			ret += FormatSpecifierStringForFormatSpecifier(segment.(model.FormatSpecifierSegment))
		}
	}
	return ret
}

func GetStringsFileContents(set model.TranslationSet, language string) string {
	ret := "/**\n" +
		" * Generated by `Sanat`\n" +
		" * Language: " + language + "\n" +
		" */\n\n"
	for _, section := range set.Sections {
		sectionHeadingPrinted := false
		for _, translation := range section.Translations {
			if !translation.IsForPlatform(model.PlatformApple) {
				continue
			}

			value := translation.ValueForLanguage(language)
			if value == nil {
				continue
			}

			if !sectionHeadingPrinted && 0 < len(section.Name) {
				ret += "\n/********** " + escapedForComment(section.Name) + " **********/\n\n"
				sectionHeadingPrinted = true
			}

			if 0 < len(translation.Comment) {
				ret += "/* " + escapedForComment(translation.Comment) + " */\n"
			}
			ret += fmt.Sprintf("\"%s\" = \"%s\";\n",
				translation.Key,
				StringFromSegments(value.Segments))
		}
	}
	return ret
}

func WriteStringsFiles(set model.TranslationSet, outDirPath string) {
	for language, _ := range set.Languages {
		lprojPath := path.Join(outDirPath, language+".lproj")
		os.MkdirAll(lprojPath, 0777)

		f, err := os.Create(path.Join(lprojPath, "Localizable.strings"))
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(GetStringsFileContents(set, language))
		if err != nil {
			panic(err)
		}
	}
}
