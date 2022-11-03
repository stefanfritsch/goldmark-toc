package toc

import (
	fences "github.com/stefanfritsch/goldmark-fences"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

const _defaultTitle = "Table of Contents"

// Transformer is a Goldmark AST transformer adds a TOC to the top of a
// Markdown document.
//
// To use this, either install the Extender on the goldmark.Markdown object,
// or install the AST transformer on the Markdown parser like so.
//
//	markdown := goldmark.New(...)
//	markdown.Parser().AddOptions(
//	  parser.WithAutoHeadingID(),
//	  parser.WithASTTransformers(
//	    util.Prioritized(&toc.Transformer{}, 100),
//	  ),
//	)
//
// NOTE: Unless you've supplied your own parser.IDs implementation, you'll
// need to enable the WithAutoHeadingID option on the parser to generate IDs
// and links for headings.
type Transformer struct {
	// Title is the title of the table of contents section.
	// Defaults to "Table of Contents" if unspecified.
	Title     string
	AddFences bool   // Wrap the toc in a div
	FencesID  string // html-id for the wrapping div if AddFences=true. Default is #md-toc
	PruneToc  bool   // if true remove wrapping lists for heading levels above the highest heading (e.g. if you don't have an h1 heading we won't render a list for h1 headings)
}

var _ parser.ASTTransformer = (*Transformer)(nil) // interface compliance

// Transform adds a table of contents to the provided Markdown document.
//
// Errors encountered while transforming are ignored. For more fine-grained
// control, use Inspect and transform the document manually.
func (t *Transformer) Transform(doc *ast.Document, reader text.Reader, pctx parser.Context) {
	toc, err := Inspect(doc, reader.Source(), &NewInspectOption{pruneToc: t.PruneToc})
	if err != nil {
		// There are currently no scenarios under which Inspect
		// returns an error but we have to account for it anyway.
		return
	}

	// Don't add anything for documents with no headings.
	if len(toc.Items) == 0 {
		return
	}

	if t.AddFences {
		var fid string

		if t.FencesID == "" {
			fid = "md-toc"
		} else {
			fid = t.FencesID
		}
		node := fences.NewFencedContainer()
		node.SetAttributeString("id", []byte(fid))
		node.SetAttributeString("class", []byte("toc nav elem-nav"))
		insertTOC(node, toc, t.Title)

		doc.InsertBefore(doc, doc.FirstChild(), node)
	} else {
		insertTOC(doc, toc, t.Title)
		// doc.InsertBefore(doc, doc.FirstChild(), tocList)
	}
}

func insertTOC(node ast.Node, toc *TOC, title string) {
	tocList := RenderList(toc)

	node.InsertBefore(node, node.FirstChild(), tocList)

	if len(title) == 0 {
		title = _defaultTitle
	}
	heading := ast.NewHeading(1)
	heading.AppendChild(heading, ast.NewString([]byte(title)))
	node.InsertBefore(node, node.FirstChild(), heading)
}
