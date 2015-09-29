package base

import (
	"hasseg.org/sanat/model"
)

// NoOpPreprocessor is a base type that does nothing. It provides the
// boilerplate for processor endpoints that "concrete" preprocessors are not
// interested in using.
type NoOpPreprocessor struct{}

func (pp NoOpPreprocessor) ProcessRawValue(s string) string {
	return s
}
func (pp NoOpPreprocessor) ProcessValueSegments(segments []model.Segment) []model.Segment {
	return segments
}
