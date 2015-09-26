package parser

import (
    "os"
    "bufio"
    "errors"
    "fmt"
    "strings"
    "strconv"
    "hasseg.org/sanat/model"
)

type translationParser struct {
    lineNumber int
    numErrors int
}

func (p *translationParser) reportError(message string) {
    p.numErrors++
    fmt.Fprintln(os.Stderr, "ERROR on line", p.lineNumber, message)
}

func (p *translationParser) intFromString(s string) int {
    parsedInt, err := strconv.Atoi(s)
    if err == nil {
        return parsedInt
    } else {
        p.reportError(err.Error())
        return 0
    }
}

func (p *translationParser) newFormatSpecifierSegmentFromSpecifierText(text string) model.TranslationValueSegment {
    s := strings.TrimRight(strings.TrimLeft(text, "{"), "}")

    // Read (potential) semantic order index
    // {a:xxx}
    //  ~
    semanticOrderIndex := -1
    orderSeparatorIndex := strings.Index(s, ":")
    if orderSeparatorIndex != -1 {
        if 0 < orderSeparatorIndex {
            semanticOrderIndex = p.intFromString(s[0:orderSeparatorIndex])
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
                numDecimals = p.intFromString(decimalCountString)
            }
        }
    }

    return model.NewFormatSpecifierSegment(dataType, numDecimals, semanticOrderIndex)
}

func componentsInCommaSeparatedList(text string) []string {
    ret := make([]string, 0)
    for _,s := range strings.Split(text, ",") {
        ret = append(ret, strings.TrimSpace(s))
    }
    return ret
}

func (p *translationParser) platformsFromCommaSeparatedString(text string) []model.TranslationPlatform {
    ret := make([]model.TranslationPlatform, 0)
    for _,s := range componentsInCommaSeparatedList(text) {
        platform := model.PlatformNone
        switch strings.ToLower(s) {
            case "apple": platform = model.PlatformApple
            case "android": platform = model.PlatformAndroid
            case "windows": platform = model.PlatformWindows
        }
        if platform == model.PlatformNone {
            p.reportError("Unknown platform value: '"+s+"'")
        } else {
            ret = append(ret, platform)
        }
    }
    return ret
}

func (p *translationParser) newSegmentsFromValue(text string) []model.TranslationValueSegment {
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
            ret = append(ret, p.newFormatSpecifierSegmentFromSpecifierText(specifierText))
            continue
        }
        accumulatedString += c
    }

    if 0 < len(accumulatedString) {
        ret = append(ret, model.NewTextSegment(accumulatedString))
    }

    return ret
}

func (p *translationParser) parseTranslationSet(inputPath string) model.TranslationSet {
    f, err := os.Open(inputPath)
    if err != nil {
        panic(err)
    }
    lineScanner := bufio.NewScanner(f)

    set := model.NewTranslationSet()
    var currentSection *model.TranslationSection
    var currentTranslation *model.Translation

    for lineScanner.Scan() {
        p.lineNumber++

        rawLine := lineScanner.Text()
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
                p.reportError("Loose line not in a translation block: " + rawLine)
            } else {
                separatorIndex := strings.Index(trimmedLine, "=")
                if separatorIndex == -1 {
                    p.reportError("Cannot find separator '=' on line: " + rawLine)
                } else {
                    key := strings.TrimSpace(trimmedLine[0:separatorIndex])
                    value := strings.TrimSpace(trimmedLine[separatorIndex+1:])
                    if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
                        value = value[1:len(value)-1]
                    }
                    if strings.ToLower(key) == "platforms" {
                        currentTranslation.Platforms = p.platformsFromCommaSeparatedString(value)
                    } else {
                        currentTranslation.AddValue(key, p.newSegmentsFromValue(value))
                        set.Languages[key] = true
                    }
                }
            }
        }
    }

    if err := lineScanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "ERROR: Parser error while reading file:", err)
    }

    return set
}

func NewTranslationSetFromFile(inputPath string) (model.TranslationSet, error) {
    parser := translationParser{}
    ret := parser.parseTranslationSet(inputPath)
    if parser.numErrors == 0 {
        return ret, nil
    } else {
        return ret, errors.New("Errors while parsing")
    }
}
