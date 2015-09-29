package smartypants

import (
	"bytes"
	"html"

	"github.com/kr/smartypants"

	"hasseg.org/sanat/preprocessing/base"
)

func textProcessedBySmartypants(text string) string {
	buffer := new(bytes.Buffer)
	smartypants.New(buffer, 0).Write([]byte(text))
	ret := buffer.String()

	// At this point `ret` contains HTML special entities like &lt;
	// so we gotta get rid of them.

	return html.UnescapeString(ret)
}

type Preprocessor struct {
	base.NoOpPreprocessor
}

func (pp Preprocessor) ProcessRawValue(v string) string {
	return textProcessedBySmartypants(v)
}
