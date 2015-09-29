package preprocessing

import (
	"errors"

	"hasseg.org/sanat/model"
	"hasseg.org/sanat/preprocessing/base"
	"hasseg.org/sanat/preprocessing/markdown"
	"hasseg.org/sanat/preprocessing/smartypants"
)

type PreProcessor interface {
	ProcessRawValue(string) string
	ProcessValueSegments([]model.Segment) []model.Segment
}

// GroupPreProcessor simply wraps a group of "concrete" preprocessors and
// invokes all of them transparently.
type GroupPreProcessor struct {
	ConcreteProcessors []PreProcessor
}

func (pp GroupPreProcessor) ProcessRawValue(s string) string {
	ret := s
	for _, processor := range pp.ConcreteProcessors {
		ret = processor.ProcessRawValue(ret)
	}
	return ret
}
func (pp GroupPreProcessor) ProcessValueSegments(segments []model.Segment) []model.Segment {
	ret := segments
	for _, processor := range pp.ConcreteProcessors {
		ret = processor.ProcessValueSegments(ret)
	}
	return ret
}

func PreProcessorForName(name string) (PreProcessor, error) {
	var PreProcessorsByName = map[string]PreProcessor{
		"markdown":    markdown.PreProcessor{},
		"smartypants": smartypants.PreProcessor{},
	}

	ret := PreProcessorsByName[name]
	if ret != nil {
		return ret, nil
	}

	e := "Unknown preprocessor '" + name + "' — available preprocessors: "
	for preprocessorName, _ := range PreProcessorsByName {
		e += preprocessorName + " "
	}
	return nil, errors.New(e)
}

func GroupPreProcessorForProcessorNames(names []string) (PreProcessor, error) {
	ret := make([]PreProcessor, 0)
	for _, name := range names {
		processor, err := PreProcessorForName(name)
		if err != nil {
			return nil, err
		}
		ret = append(ret, processor)
	}
	return GroupPreProcessor{ConcreteProcessors: ret}, nil
}

func NewNoOpPreProcessor() PreProcessor {
	return base.NoOpPreProcessor{}
}
