package java

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/util"
)

func FormatSpecifierStringForFormatSpecifier(segment model.FormatSpecifierSegment, index int) string {
	outputOrderIndex := index
	if 0 < segment.SemanticOrderIndex {
		outputOrderIndex = segment.SemanticOrderIndex - 1
	}
	ret := "{" + strconv.Itoa(outputOrderIndex)

	if segment.DataType == model.DataTypeInteger {
		ret += ",number,integer"
	} else if segment.DataType == model.DataTypeFloat {
		ret += ",number"
		if 0 == segment.NumberOfDecimals {
			ret += ",#"
		} else if 0 < segment.NumberOfDecimals {
			ret += ",#." + strings.Repeat("#", segment.NumberOfDecimals)
		}
	}

	ret += "}"
	return ret
}

var _moustacheRegexp *regexp.Regexp = nil

func SanitizedForStringValue(text string) string {
	// The escaping/quoting rules for Java MessageFormat strings are a
	// MASSIVE PAIN IN THE ASS. Let’s just do our best.
	// https://docs.oracle.com/javase/7/docs/api/java/text/MessageFormat.html
	//
	if _moustacheRegexp == nil {
		_moustacheRegexp, _ = regexp.Compile("([{}]+)")
	}
	s := ""
	for index, part := range strings.Split(text, "'") {
		p := _moustacheRegexp.ReplaceAllString(part, "'$1'")
		if index == 0 || strings.HasSuffix(s, "'") && strings.HasPrefix(p, "'") {
			s += p
		} else {
			s += "''" + p
		}
	}
	return util.XMLEscaped(s)
}

func SanitizedForKey(text string) string {
	return util.XMLEscaped(text)
}

func stringFromSegments(segments []model.Segment) string {
	ret := ""
	for index, segment := range segments {
		switch segment.(type) {
		case model.TextSegment:
			ret += SanitizedForStringValue(segment.(model.TextSegment).Text)
		case model.FormatSpecifierSegment:
			ret += FormatSpecifierStringForFormatSpecifier(segment.(model.FormatSpecifierSegment), index)
		}
	}
	return ret
}

// JDK 5 adds the API `Properties.loadFromXML(InputStream)`.
// We're targeting that instead of the "classic" .properties files,
// which are ISO-8859-1 encoded.
//
func GetPropertiesFileContents(set model.TranslationSet, language string) string {
	ret := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n" +
		"<!DOCTYPE properties SYSTEM \"http://java.sun.com/dtd/properties.dtd\">\n" +
		"<!--\n" +
		"Java XML Properties File\n" +
		"Generated by Sanat\n" +
		"Language: " + language + "\n" +
		"-->\n" +
		"<properties>\n"

	for _, section := range set.Sections {
		sectionHeadingPrinted := false
		for _, translation := range section.Translations {
			if !translation.IsForPlatform(model.PlatformJava) {
				continue
			}

			value := translation.ValueForLanguage(language)
			if value == nil {
				continue
			}

			if !sectionHeadingPrinted && 0 < len(section.Name) {
				sanitizedSectionName := strings.Replace(section.Name, "--", "- -", -1)
				ret += "\n  <!-- ********** " + sanitizedSectionName + " ********** -->\n\n"
				sectionHeadingPrinted = true
			}

			ret += fmt.Sprintf("  <entry key=\"%s\">%s</entry>\n", SanitizedForKey(translation.Key), stringFromSegments(value.Segments))
		}
	}
	ret += "</properties>\n"
	return ret
}

func WritePropertiesFiles(set model.TranslationSet, outDirPath string) {
	for language, _ := range set.Languages {
		os.MkdirAll(outDirPath, 0777)

		f, err := os.Create(path.Join(outDirPath, "Properties_"+language+".xml"))
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(GetPropertiesFileContents(set, language))
		if err != nil {
			panic(err)
		}
	}
}
