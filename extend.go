package toc

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Extender extends a Goldmark Markdown parser and renderer to always include
// a table of contents in the output.
//
// To use this, install it into your Goldmark Markdown object.
//
//	md := goldmark.New(
//	  // ...
//	  goldmark.WithParserOptions(parser.WithAutoHeadingID()),
//	  goldmark.WithExtensions(
//	    // ...
//	    &toc.Extender{
//	    },
//	  ),
//	)
//
// This will install the default Transformer. For more control, install the
// Transformer directly on the Markdown Parser.
//
// NOTE: Unless you've supplied your own parser.IDs implementation, you'll
// need to enable the WithAutoHeadingID option on the parser to generate IDs
// and links for headings.
type Extender struct {
	Title     string
	AddFences bool
	FencesID  string
	PruneTOC  bool
}

// Extend adds support for rendering a table of contents to the provided
// Markdown parser/renderer.
func (e *Extender) Extend(md goldmark.Markdown) {
	md.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&Transformer{Title: e.Title, AddFences: e.AddFences, FencesID: e.FencesID, PruneTOC: e.PruneTOC}, 100),
		),
	)
}
