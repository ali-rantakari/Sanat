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

    // Read data type indicator
    dataType := model.DataTypeObject
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

        line := strings.TrimSpace(scanner.Text())

        if len(line) == 0 {
            continue
        }

        if strings.HasPrefix(line, "[[") && strings.HasSuffix(line, "]]") {
            currentSection = set.AddSection(line[2:len(line)-2])
        } else if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
            if currentSection == nil {
                ReportParserError(lineNumber, "Loose translation not in a section: " + line)
            } else {
                currentTranslation = currentSection.AddTranslation(line[1:len(line)-1])
            }
        } else {
            if currentTranslation == nil {
                ReportParserError(lineNumber, "Loose line not in a translation block: " + line)
            } else {
                separatorIndex := strings.Index(line, "=")
                if separatorIndex == -1 {
                    ReportParserError(lineNumber, "Cannot find separator '=' on line: " + line)
                } else {
                    key := strings.TrimSpace(line[0:separatorIndex])
                    value := strings.TrimSpace(line[separatorIndex+1:])
                    currentTranslation.AddValue(key, NewSegmentsFromValue(value))
                }
            }
        }
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading file:", err)
    }

    return set
}
