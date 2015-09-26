package output

import (
    "fmt"
    "strconv"
    "hasseg.org/sanat/model"
)

func StringForFormatSpecifier(segment model.TranslationValueSegment) string {
    ret := ""
    switch segment.DataType {
        case model.DataTypeString: ret += "<string>"
        case model.DataTypeInteger: ret += "<integer>"
        case model.DataTypeFloat: ret += "<float>"
        case model.DataTypeObject: ret += "<object>"
        default: ret += "<??>"
    }
    if segment.DataType == model.DataTypeFloat && 0 <= segment.NumberOfDecimals {
        ret += ", " + strconv.Itoa(segment.NumberOfDecimals) + " decimals"
    }
    if 0 < segment.SemanticOrderIndex {
        ret += ", order #" + strconv.Itoa(segment.SemanticOrderIndex)
    }
    return ret
}

func StringForPlatform(platform model.TranslationPlatform) string {
    switch platform {
        case model.PlatformApple: return "Apple"
    }
    return "??"
}

func DumpTranslationSet(set model.TranslationSet, outputDirPath string) {
    fmt.Println("Languages:", set.Languages)
    for _,section := range set.Sections {
        fmt.Println("Section: " + section.Name)
        for _,translation := range section.Translations {
            fmt.Println("  Translation: " + translation.Key)
            if 0 < len(translation.Platforms) {
                for _,platform := range translation.Platforms {
                    fmt.Println("    Platform: " + StringForPlatform(platform))
                }
            }
            for _,value := range translation.Values {
                fmt.Println("    Language: " + value.Language)
                for _,segment := range value.Segments {
                    if segment.IsFormatSpecifier {
                        fmt.Println("       fmt: " + StringForFormatSpecifier(segment))
                    } else {
                        fmt.Println("      Text: '" + segment.Text + "'")
                    }
                }
            }
        }
    }
}
