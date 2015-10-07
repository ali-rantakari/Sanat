package parser

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/preprocessing"
	"hasseg.org/sanat/util"
)

type ParserErrorHandler func(lineNumber int, message string)

type translationParser struct {
	lineNumber   int
	numErrors    int
	errorHandler ParserErrorHandler
}

func (p *translationParser) reportError(message string) {
	p.numErrors++
	if p.errorHandler != nil {
		p.errorHandler(p.lineNumber, message)
	}
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

func (p *translationParser) formatSpecifierSegmentFromSpecifierText(text string) model.Segment {
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
	case "@":
		dataType = model.DataTypeObject
	case "f":
		dataType = model.DataTypeFloat
	case "d":
		dataType = model.DataTypeInteger
	case "s":
		dataType = model.DataTypeString
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

func (p *translationParser) platformsFromCommaSeparatedString(text string) []model.TranslationPlatform {
	ret := make([]model.TranslationPlatform, 0)
	for _, s := range util.ComponentsFromCommaSeparatedList(text) {
		platform := model.PlatformNone
		switch strings.ToLower(s) {
		case "apple":
			platform = model.PlatformApple
		case "android":
			platform = model.PlatformAndroid
		case "windows":
			platform = model.PlatformWindows
		}
		if platform == model.PlatformNone {
			p.reportError("Unknown platform value: '" + s + "'")
		} else {
			ret = append(ret, platform)
		}
	}
	return ret
}

func (p *translationParser) segmentsFromTranslationValueString(text string) []model.Segment {
	ret := make([]model.Segment, 0)

	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanRunes)

	scanUntilEndOfFormatSpecifier := func(scanner *bufio.Scanner) string {
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
			ret = append(ret, p.formatSpecifierSegmentFromSpecifierText(specifierText))
			continue
		}
		accumulatedString += c
	}

	if 0 < len(accumulatedString) {
		ret = append(ret, model.NewTextSegment(accumulatedString))
	}

	return ret
}

func (p *translationParser) parseTranslationSet(inputReader io.Reader, preprocessor preprocessing.Preprocessor) model.TranslationSet {
	lineScanner := bufio.NewScanner(inputReader)

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

		processSectionHeadingRow := func() {
			if strings.HasPrefix(rawLine, "===") {
				currentSection = set.AddSection(strings.Trim(trimmedLine, "= "))
			} else {
				p.reportError("Unknown un-indented line '" + rawLine + "' — Prepend with === if section; indent if translation key.")
			}
		}

		processTranslationKeyHeadingRow := func() {
			if currentSection == nil { // Add implicit default section if needed
				currentSection = set.AddSection("")
			}
			if currentTranslation != nil {
				if len(currentTranslation.Values) == 0 {
					p.reportError("Translation '" + currentTranslation.Key + "' has no values")
				}
			}
			currentTranslation = currentSection.AddTranslation(trimmedLine)
		}

		processTranslationMetadataRow := func() {
			if currentTranslation == nil {
				p.reportError("Loose line not in a translation block: " + rawLine)
				return
			}

			separatorIndex := strings.Index(trimmedLine, "=")
			if separatorIndex == -1 {
				p.reportError("Cannot find separator '=' on line: " + rawLine)
				return
			}

			key := strings.TrimSpace(trimmedLine[0:separatorIndex])
			value := strings.TrimSpace(trimmedLine[separatorIndex+1:])
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = value[1 : len(value)-1]
			}

			lowerKey := strings.ToLower(key)
			if lowerKey == "platforms" {
				currentTranslation.Platforms = p.platformsFromCommaSeparatedString(value)
			} else if lowerKey == "tags" {
				currentTranslation.Tags = util.ComponentsFromCommaSeparatedList(value)
			} else if lowerKey == "comment" {
				currentTranslation.Comment = value
			} else {
				value = preprocessor.ProcessRawValue(value)
				segments := preprocessor.ProcessValueSegments(p.segmentsFromTranslationValueString(value))
				currentTranslation.AddValue(key, segments)
				set.Languages[key] = true
			}
		}

		leadingWhitespaceCount := len(util.LeadingWhitespace(rawLine))

		if leadingWhitespaceCount == 0 {
			processSectionHeadingRow()
		} else if leadingWhitespaceCount == 2 {
			processTranslationKeyHeadingRow()
		} else if leadingWhitespaceCount == 4 {
			processTranslationMetadataRow()
		}
	}

	if err := lineScanner.Err(); err != nil {
		p.reportError("Error while reading file: " + err.Error())
	}

	return set
}

func TranslationSetFromFile(inputPath string, preprocessor preprocessing.Preprocessor, errorHandler ParserErrorHandler) (model.TranslationSet, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	parser := translationParser{errorHandler: errorHandler}
	ret := parser.parseTranslationSet(f, preprocessor)
	if parser.numErrors == 0 {
		return ret, nil
	} else {
		return ret, errors.New("Errors while parsing")
	}
}
