package preprocessing

import (
	"strings"

	"github.com/russross/blackfriday"
	"hasseg.org/sanat/model"

	"hasseg.org/sanat/util"
)

func htmlFromMarkdown(md string) string {
	if len(md) == 0 {
		return md
	}

	// Render markdown to HTML
	//
	htmlFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	markdownExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS

	ret := string(blackfriday.MarkdownOptions(
		[]byte(md),
		blackfriday.HtmlRenderer(htmlFlags, "", ""),
		blackfriday.Options{Extensions: markdownExtensions}))

	// Remove wrapping paragraph tags
	//
	ret = strings.TrimSpace(ret)
	if strings.HasPrefix(ret, "<p>") && strings.HasSuffix(ret, "</p>") {
		ret = ret[3 : len(ret)-4]
	}

	// Reinstate possible leading/trailing whitespace
	// (the markdown compiler will have removed it)
	//
	if leadingSpaces := util.LeadingWhitespace(md); 0 < len(leadingSpaces) {
		ret = leadingSpaces + ret
	}
	if trailingSpaces := util.TrailingWhitespace(md); 0 < len(trailingSpaces) {
		ret = ret + trailingSpaces
	}

	return ret
}

func PreprocessMarkdown(set *model.TranslationSet) error {
	//
	// TODO: need to compile the translation string as a whole — compiling
	// each segment in isolation will yield incorrect results.
	// e.g.:
	//
	//   _Eka {0} toka_
	//
	// should yield:
	//
	//   <em>Eka {0} toka</em>
	//
	set.IterateTranslationValues(func(value *model.TranslationValue) {
		for g := 0; g < len(value.Segments); g++ {
			if !value.Segments[g].IsFormatSpecifier {
				value.Segments[g].Text = htmlFromMarkdown(value.Segments[g].Text)
			}
		}
	})
	return nil
}
