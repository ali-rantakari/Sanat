package util

import (
	"bytes"
	"encoding/xml"
	"io"
)

func XMLEscaped(s string) string {
	var b bytes.Buffer
	xml.EscapeText(&b, []byte(s))
	return b.String()
}

func XMLIsValid(xmlString string) bool {
	decoder := xml.NewDecoder(bytes.NewBufferString(xmlString))
	for {
		_, err := decoder.Token()
		if err == nil {
			continue
		} else if err == io.EOF {
			break
		}
		return false
	}
	return true
}
