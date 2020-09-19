package models

import (
	"bytes"
	"html/template"

	bfp "github.com/Depado/bfplus"
	"github.com/alecthomas/chroma/formatters/html"
	bf "gopkg.in/russross/blackfriday.v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var exts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.DefinitionLists | bf.Footnotes

var flags = bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.SmartypantsLatexDashes | bf.TOC

// GlobCSS is a byte slice containing the style CSS of the renderer
var GlobCSS template.CSS

func render(input []byte) []byte {
	r := bfp.NewRenderer(
		bfp.WithAdmonition(),
		bfp.WithCodeHighlighting(
			bfp.Style(viper.GetString("blog.code.style")),
			bfp.WithoutAutodetect(),
			bfp.ChromaOptions(html.WithClasses(true)),
		),
		bfp.WithHeadingAnchors(),
		bfp.Extend(
			bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: flags}),
		),
	)

	if GlobCSS == "" && r.Highlighter.Formatter.Classes {
		b := new(bytes.Buffer)
		if err := r.Highlighter.Formatter.WriteCSS(b, r.Highlighter.Style); err != nil {
			logrus.WithError(err).Warning("Couldn't write CSS")
		}
		GlobCSS = template.CSS(b.String())
	}

	return bf.Run(
		input,
		bf.WithRenderer(r),
		bf.WithExtensions(exts),
	)
}
