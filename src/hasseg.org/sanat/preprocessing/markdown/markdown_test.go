package markdown_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"hasseg.org/sanat/preprocessing/markdown"
)

func TestMarkdownToHTML(t *testing.T) {
	pp := markdown.Preprocessor{}

	ass := func(expected string, input string) {
		assert.Equal(t, expected, pp.ProcessRawValue(input), input)
	}

	ass("", "")
	ass("Hello", "Hello")

	// Basic Markdown on individual word
	ass("<em>Hello</em>", "_Hello_")
	ass("<em>Hello</em>", "*Hello*")
	ass("<strong>Hello</strong>", "__Hello__")
	ass("<strong>Hello</strong>", "**Hello**")
	ass("<code>Hello</code>", "`Hello`")

	// Preserves leading/trailing whitespace
	ass("\t<em>Hello</em>\t", "\t_Hello_\t")
	ass("  <em>Hello</em>  ", "  _Hello_  ")
}
