package output

import (
    "fmt"
    "strconv"
    "hasseg.org/sanat/model"
)

func AppleFormatSpecifierStringForFormatSpecifier(segment model.TranslationValueSegment) string {
    ret := "%"
    if 0 < segment.SemanticOrderIndex {
        ret += strconv.Itoa(segment.SemanticOrderIndex) + "$"
    }
    if segment.DataType == model.DataTypeFloat && 0 <= segment.NumberOfDecimals {
        ret += "." + strconv.Itoa(segment.NumberOfDecimals)
    }
    switch segment.DataType {
        case model.DataTypeString: ret += "s"
        case model.DataTypeInteger: ret += "d"
        case model.DataTypeFloat: ret += "f"
        case model.DataTypeObject: ret += "@"
    }
    return ret
}

func AppleStringFromSegments(segments []model.TranslationValueSegment) string {
    ret := ""
    for _,segment := range segments {
        if segment.IsFormatSpecifier {
            ret += AppleFormatSpecifierStringForFormatSpecifier(segment)
        } else {
            ret += segment.Text
        }
    }
    return ret
}

func WriteAppleStringsFile(set model.TranslationSet, language string) {
    fmt.Println("/**\n" +
                " * Apple Strings File\n" +
                " * Generated by TODO\n" +
                " * Language: " + language + "\n" +
                " */")
    for _,section := range set.Sections {
        fmt.Println("\n/********** " + section.Name + " **********/\n")
        for _,translation := range section.Translations {
            for _,value := range translation.Values {
                if value.Language == language {
                    fmt.Printf("\"%s\" = \"%s\";\n",
                               translation.Key,
                               AppleStringFromSegments(value.Segments))
                }
            }
        }
    }
}
