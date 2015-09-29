package android

import (
	"bytes"
	"encoding/xml"
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
	case model.DataTypeObject:
		fallthrough
	case model.DataTypeString:
		ret += "s"
	case model.DataTypeInteger:
		ret += "d"
	case model.DataTypeFloat:
		ret += "f"
	}
	return ret
}

func xmlEscaped(text string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(text))
	return b.String()
}

func SanitizedForString(text string) string {
	return xmlEscaped(strings.Replace(text, "%", "%%", -1))
}

func stringFromSegments(segments []model.Segment) string {
	ret := ""
	for _, segment := range segments {
		switch segment.(type) {
		case model.TextSegment:
			ret += SanitizedForString(segment.(model.TextSegment).Text)
		case model.FormatSpecifierSegment:
			ret += FormatSpecifierStringForFormatSpecifier(segment.(model.FormatSpecifierSegment))
		}
	}
	return ret
}

func sanitizedForXMLComment(s string) string {
	return strings.Replace(s, "--", "- -", -1)
}

func GetStringsFileContents(set model.TranslationSet, language string) string {
	ret := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n<resources>\n"
	for _, section := range set.Sections {
		if 0 < len(section.Name) {
			ret += "\n    <!-- ********** " + sanitizedForXMLComment(section.Name) + " ********** -->\n\n"
		}
		for _, translation := range section.Translations {
			if !translation.IsForPlatform(model.PlatformAndroid) {
				continue
			}
			for _, value := range translation.Values {
				if value.Language == language {
					if 0 < len(translation.Comment) {
						ret += fmt.Sprintf("    <!-- %s -->\n",
							sanitizedForXMLComment(translation.Comment))
					}
					ret += fmt.Sprintf("    <string name=\"%s\">%s</string>\n",
						xmlEscaped(translation.Key),
						stringFromSegments(value.Segments))
				}
			}
		}
	}
	ret += "</resources>\n"
	return ret
}

func WriteStringsFiles(set model.TranslationSet, outDirPath string) {
	for language, _ := range set.Languages {
		valuesDirPath := path.Join(outDirPath, "values-"+language)
		os.MkdirAll(valuesDirPath, 0777)

		f, err := os.Create(path.Join(valuesDirPath, "strings.xml"))
		if err != nil {
			panic(err)
		}

		_, err = f.WriteString(GetStringsFileContents(set, language))
		if err != nil {
			panic(err)
		}
	}
}
