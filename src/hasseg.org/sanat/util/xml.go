package util

import (
	"bytes"
	"encoding/xml"
)

func XMLEscaped(s string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(s))
	return b.String()
}
