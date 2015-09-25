package output

import (
    "fmt"
    "hasseg.org/sanat/model"
)

func DumpTranslationSet(set model.TranslationSet) {
    for _,section := range set.Sections {
        fmt.Println("Section: " + section.Name)
        for _,translation := range section.Translations {
            fmt.Println("  Translation: " + translation.Key)
            for _,value := range translation.Values {
                fmt.Println("    Value: " + value.Language + " = " + value.Text)
            }
        }
    }
}
