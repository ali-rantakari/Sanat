package parser

import (
    "os"
    "bufio"
    "fmt"
    "strings"
    "strconv"
    "hasseg.org/sanat/model"
)

func ReportParserError(lineNumber int, message string) {
    fmt.Fprintln(os.Stderr, "ERROR: Parser error:", message)
}

func IntFromString(s string) int {
    parsedInt, err := strconv.Atoi(s)
    if err == nil {
        return parsedInt
    } else {
        ReportParserError(0, err.Error())
        return 0
    }
}

func NewFormatSpecifierSegmentFromSpecifierText(text string) model.TranslationValueSegment {
    s := strings.TrimRight(strings.TrimLeft(text, "{"), "}")

    // Read (potential) semantic order index
    // {a:xxx}
    //  ~
    semanticOrderIndex := -1
    orderSeparatorIndex := strings.Index(s, ":")
    if orderSeparatorIndex != -1 {
        if 0 < orderSeparatorIndex {
            semanticOrderIndex = IntFromString(s[0:orderSeparatorIndex])
        }
        s = s[orderSeparatorIndex+1:]
    }
    if semanticOrderIndex == 0 { // normalize: -1 means "none"
        semanticOrderIndex = -1
    }

    defaultDataType := model.DataTypeObject

    if len(s) == 0 {
        return model.NewFormatSpecifierSegment(defaultDataType, -1, -1)
    }

    // Read data type indicator
    dataType := defaultDataType
    switch strings.ToLower(s[0:1]) {
        case "@": dataType = model.DataTypeObject
        case "f": dataType = model.DataTypeFloat
        case "d": dataType = model.DataTypeInteger
        case "s": dataType = model.DataTypeString
    }
    s = s[1:]

    // Read (potential) decimal count
    numDecimals := -1
    if dataType == model.DataTypeFloat {
        decimalCountIndex := strings.Index(s, ".")
        if decimalCountIndex != -1 {
            decimalCountString := s[decimalCountIndex+1:]
            if 0 < len(decimalCountString) {
                numDecimals = IntFromString(decimalCountString)
            }
        }
    }

    return model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex)
}

func NewSegmentsFromValue(text string) []model.TranslationValueSegment {
    ret := make([]model.TranslationValueSegment, 0)

    scanner := bufio.NewScanner(strings.NewReader(text))
    scanner.Split(bufio.ScanRunes)

    scanUntilEndOfFormatSpecifier := func (scanner *bufio.Scanner) string {
        accumulatedString := ""
        for scanner.Scan() {
            c := scanner.Text()
            if c == "}" {
                break
            }
            accumulatedString += c
        }
        return accumulatedString
    }

    accumulatedString := ""
    for scanner.Scan() {
        c := scanner.Text()
        if c == "\\" {
            scanner.Scan()
            accumulatedString += scanner.Text()
            continue
        }
        if c == "{" {
            ret = append(ret, model.NewTextSegment(accumulatedString))
            accumulatedString = ""

            specifierText := scanUntilEndOfFormatSpecifier(scanner)
            ret = append(ret, NewFormatSpecifierSegmentFromSpecifierText(specifierText))
            continue
        }
        accumulatedString += c
    }

    if 0 < len(accumulatedString) {
        ret = append(ret, model.NewTextSegment(accumulatedString))
    }

    return ret
}

func NewTranslationSetFromFile(inputPath string) model.TranslationSet {
    f, err := os.Open(inputPath)
    if err != nil {
        panic(err)
    }
    scanner := bufio.NewScanner(f)

    set := model.NewTranslationSet()
    var currentSection *model.TranslationSection
    var currentTranslation *model.Translation

    lineNumber := 0
    for scanner.Scan() {
        lineNumber++

        rawLine := scanner.Text()
        trimmedLine := strings.TrimSpace(rawLine)

        if len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#") {
            continue
        }

        if strings.HasPrefix(trimmedLine, "===") {
            currentSection = set.AddSection(strings.Trim(trimmedLine, "= "))
        } else if !strings.HasPrefix(rawLine, "  ") {
            if currentSection == nil { // Add implicit default section if needed
                currentSection = set.AddSection("")
            }
            currentTranslation = currentSection.AddTranslation(trimmedLine)
        } else {
            if currentTranslation == nil {
                ReportParserError(lineNumber, "Loose line not in a translation block: " + rawLine)
            } else {
                separatorIndex := strings.Index(trimmedLine, "=")
                if separatorIndex == -1 {
                    ReportParserError(lineNumber, "Cannot find separator '=' on line: " + rawLine)
                } else {
                    language := strings.TrimSpace(trimmedLine[0:separatorIndex])
                    value := strings.TrimSpace(trimmedLine[separatorIndex+1:])
                    if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
                        value = value[1:len(value)-1]
                    }
                    currentTranslation.AddValue(language, NewSegmentsFromValue(value))
                    set.Languages[language] = true
                }
            }
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading file:", err)
    }

    return set
}
