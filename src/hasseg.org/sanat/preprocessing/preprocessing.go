package preprocessing

import (
	"errors"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/preprocessing/base"
	"hasseg.org/sanat/preprocessing/markdown"
	"hasseg.org/sanat/preprocessing/smartypants"
)

type Preprocessor interface {
	ProcessRawValue(string) string
	ProcessValueSegments([]model.Segment) []model.Segment
}

// GroupPreprocessor simply wraps a group of "concrete" Preprocessors and
// invokes all of them transparently.
type GroupPreprocessor struct {
	ConcreteProcessors []Preprocessor
}

func (pp GroupPreprocessor) ProcessRawValue(s string) string {
	ret := s
	for _, processor := range pp.ConcreteProcessors {
		ret = processor.ProcessRawValue(ret)
	}
	return ret
}
func (pp GroupPreprocessor) ProcessValueSegments(segments []model.Segment) []model.Segment {
	ret := segments
	for _, processor := range pp.ConcreteProcessors {
		ret = processor.ProcessValueSegments(ret)
	}
	return ret
}

func PreprocessorForName(name string) (Preprocessor, error) {
	var PreprocessorsByName = map[string]Preprocessor{
		"markdown":    markdown.Preprocessor{},
		"smartypants": smartypants.Preprocessor{},
	}

	ret := PreprocessorsByName[name]
	if ret != nil {
		return ret, nil
	}

	e := "Unknown preprocessor '" + name + "' — available preprocessors: "
	for PreprocessorName, _ := range PreprocessorsByName {
		e += PreprocessorName + " "
	}
	return nil, errors.New(e)
}

func GroupPreprocessorForProcessorNames(names []string) (Preprocessor, error) {
	ret := make([]Preprocessor, 0)
	for _, name := range names {
		processor, err := PreprocessorForName(name)
		if err != nil {
			return nil, err
		}
		ret = append(ret, processor)
	}
	return GroupPreprocessor{ConcreteProcessors: ret}, nil
}

func NewNoOpPreprocessor() Preprocessor {
	return base.NoOpPreprocessor{}
}
