package models

import (
	"github.com/Depado/bfchroma"
	bf "github.com/russross/blackfriday"
	"github.com/spf13/viper"
)

var exts = bf.NoIntraEmphasis | bf.Tables | bf.FencedCode | bf.Autolink |
	bf.Strikethrough | bf.SpaceHeadings | bf.BackslashLineBreak |
	bf.DefinitionLists | bf.Footnotes

var flags = bf.UseXHTML | bf.Smartypants | bf.SmartypantsFractions |
	bf.SmartypantsDashes | bf.SmartypantsLatexDashes | bf.TOC

func render(input []byte) []byte {
	return bf.Run(
		input,
		bf.WithRenderer(
			bfchroma.NewRenderer(
				bfchroma.WithoutAutodetect(),
				bfchroma.Extend(
					bf.NewHTMLRenderer(bf.HTMLRendererParameters{Flags: flags}),
				),
				bfchroma.Style(viper.GetString("blog.code.style")),
			),
		),
		bf.WithExtensions(exts),
	)
}
