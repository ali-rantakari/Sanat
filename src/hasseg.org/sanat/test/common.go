package test

import (
	"hasseg.org/sanat/model"
	"hasseg.org/sanat/parser"
	"hasseg.org/sanat/preprocessing"
)

func GetComprehensiveTestInputTranslationSet() model.TranslationSet {
	ret, _ := parser.TranslationSetFromFile("../testdata/comprehensive.sanat", preprocessing.NewNoOpPreprocessor(), nil)
	return ret
}
