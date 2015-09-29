package json

import (
	"fmt"
	"strconv"
	"strings"

	"hasseg.org/sanat/model"
)

func JSONForFormatSpecifier(segment model.FormatSpecifierSegment) string {
	ret := `{"dataType": "`
	switch segment.DataType {
	case model.DataTypeString:
		ret += "string"
	case model.DataTypeInteger:
		ret += "integer"
	case model.DataTypeFloat:
		ret += "float"
	case model.DataTypeObject:
		ret += "object"
	}
	ret += `"`
	if segment.DataType == model.DataTypeFloat && 0 <= segment.NumberOfDecimals {
		ret += `, "numberOfDecimals": ` + strconv.Itoa(segment.NumberOfDecimals)
	}
	if 0 < segment.SemanticOrderIndex {
		ret += `, "orderIndex": ` + strconv.Itoa(segment.SemanticOrderIndex)
	}
	return ret + "}"
}

func StringForPlatform(platform model.TranslationPlatform) string {
	switch platform {
	case model.PlatformApple:
		return "Apple"
	case model.PlatformAndroid:
		return "Android"
	case model.PlatformWindows:
		return "Windows"
	}
	return "??"
}

func escapedForJSON(s string) string {
	return strings.Replace(s, `"`, `\"`, -1)
}

func jsonForStringList(list []string) string {
	ret := "["
	for i, item := range list {
		if 0 < i {
			ret += ", "
		}
		ret += `"` + escapedForJSON(item) + `"`
	}
	return ret + "]"
}

func DumpTranslationSet(set model.TranslationSet, outputDirPath string) {
	languages := make([]string, 0, len(set.Languages))
	for k := range set.Languages {
		languages = append(languages, k)
	}

	ret := `{"languages": ` + jsonForStringList(languages) +
		`, "sections": [`

	for sectionIndex, section := range set.Sections {
		if 0 < sectionIndex {
			ret += ","
		}
		ret += `{"name": "` + escapedForJSON(section.Name) +
			`", "translations": [`

		for translationIndex, translation := range section.Translations {
			if 0 < translationIndex {
				ret += ","
			}
			ret += `{"key": "` + escapedForJSON(translation.Key) + `"`

			if 0 < len(translation.Platforms) {
				ret += `, "platforms": [`
				for platformIndex, platform := range translation.Platforms {
					if 0 < platformIndex {
						ret += ","
					}
					ret += `"` + escapedForJSON(StringForPlatform(platform)) + `"`
				}
				ret += `]`
			}

			ret += `, "values": [`
			for valueIndex, value := range translation.Values {
				if 0 < valueIndex {
					ret += ","
				}
				ret += `{"language": "` + escapedForJSON(value.Language) +
					`", "segments": [`

				for segmentIndex, segment := range value.Segments {
					if 0 < segmentIndex {
						ret += ","
					}

					switch segment.(type) {
					case model.TextSegment:
						ret += `{"text": "` + escapedForJSON(segment.(model.TextSegment).Text) + `"}`
					case model.FormatSpecifierSegment:
						ret += JSONForFormatSpecifier(segment.(model.FormatSpecifierSegment))
					}
				}
				ret += "]}"
			}
			ret += "]}"
		}
		ret += "]}"
	}

	ret += `]}`
	fmt.Print(ret)
}
