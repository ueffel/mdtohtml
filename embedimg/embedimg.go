package embedimg

import (
	"encoding/base64"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// HTMLRenderer is a renderer.NodeRenderer implementation that
// renders images as embedded base64.
type HTMLRenderer struct {
	html.Config
}

// NewHTMLRenderer returns a new embedimg.HTMLRenderer.
func NewHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &HTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *HTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindImage, r.renderImage)
}

func (r *HTMLRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString(`<img src="`)
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		path := string(n.Destination)
		if _, err := os.Stat(path); err == nil || !os.IsNotExist(err) {
			_, _ = w.WriteString("data:")
			if by, err := os.ReadFile(path); err == nil {
				mimeType := http.DetectContentType(by)
				_, _ = w.WriteString(mimeType)
				_, _ = w.WriteString(";base64,")
				enc := base64.NewEncoder(base64.RawStdEncoding, w)
				_, _ = enc.Write(by)
			}
		} else {
			_, _ = w.Write(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		}
	}
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(util.EscapeHTML(n.Text(source)))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

type embedImg struct{}

// EmbedImg is an extension that reads the images from img if they are local
// files and writes them as base64 in the src attribute
var EmbedImg = &embedImg{}

func (e *embedImg) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(), 500),
	))
}

// Interface guards
var (
	_ renderer.NodeRenderer = (*HTMLRenderer)(nil)
	_ goldmark.Extender     = (*embedImg)(nil)
)
