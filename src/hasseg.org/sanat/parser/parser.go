package parser

import (
    "os"
    "bufio"
    "fmt"
    "strings"
    "hasseg.org/sanat/model"
)

func ReportParserError(lineNumber int, message string) {
    fmt.Fprintln(os.Stderr, "ERROR: Parser error:", message)
}

func NewFormatSpecifierSegmentFromSpecifierText(text string) model.TranslationValueSegment {
    numDecimals := -1
    semanticOrderIndex := -1
    return model.NewFormatSpecifierSegment(model.DataTypeObject, numDecimals, semanticOrderIndex)
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
