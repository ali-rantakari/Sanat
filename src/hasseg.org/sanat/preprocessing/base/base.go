package base

import (
	"hasseg.org/sanat/model"
)

// NoOpPreProcessor is a base type that does nothing. It provides the
// boilerplate for processor endpoints that "concrete" preprocessors are not
// interested in using.
type NoOpPreProcessor struct{}

func (pp NoOpPreProcessor) ProcessRawValue(s string) string {
	return s
}
func (pp NoOpPreProcessor) ProcessValueSegments(segments []model.Segment) []model.Segment {
	return segments
}
