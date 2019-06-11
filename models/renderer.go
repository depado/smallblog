package models

import (
	"bytes"
	"html/template"

	"github.com/Depado/bfadmonition"
	"github.com/Depado/bfchroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	bf "gopkg.in/russross/blackfriday.v2"
)

var exts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.DefinitionLists | bf.Footnotes

var flags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.SmartypantsLatexDashes | bf.TOC

// GlobCSS is a byte slice containing the style CSS of the renderer
var GlobCSS template.CSS

func render(input []byte) []byte {
	r := bfchroma.NewRenderer(
		bfchroma.WithoutAutodetect(),
		bfchroma.Extend(
			bfadmonition.NewRenderer(
				bfadmonition.Extend(
					bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: flags}),
				),
			),
		),

		bfchroma.Style(viper.GetString("blog.code.style")),
		bfchroma.ChromaOptions(html.WithClasses()),
	)
	if GlobCSS == "" && r.Formatter.Classes {
		b := new(bytes.Buffer)
		if err := r.Formatter.WriteCSS(b, r.Style); err != nil {
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
