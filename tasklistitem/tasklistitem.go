package tasklistitem

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// HTMLRenderer is a renderer.NodeRenderer implementation that renders list
// items with CSS class "task-list-item" if it contains a TaskCheckBox
type HTMLRenderer struct {
	html.Config
}

// NewHTMLRenderer returns a new tasklist.HTMLRenderer.
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
	reg.Register(ast.KindListItem, r.renderListItem)
}

// IsTasklistItem returns true, if a ListItem contains a TaskCheckBox
func (r *HTMLRenderer) IsTasklistItem(node ast.Node) bool {
	txtBlock := node.FirstChild()
	if txtBlock == nil {
		return false
	}
	chkbox := txtBlock.FirstChild()
	if chkbox == nil {
		return false
	}
	_, ok := chkbox.(*east.TaskCheckBox)
	return ok
}

func (r *HTMLRenderer) renderListItem(w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		_, _ = w.WriteString("<li")
		if r.IsTasklistItem(n) {
			_, _ = w.WriteString(` class="task-list-item"`)
		}
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, html.ListItemAttributeFilter)
		}
		_ = w.WriteByte('>')

		fc := n.FirstChild()
		if fc != nil {
			if _, ok := fc.(*ast.TextBlock); !ok {
				_ = w.WriteByte('\n')
			}
		}
	} else {
		_, _ = w.WriteString("</li>\n")
	}
	return ast.WalkContinue, nil
}

type taskListItemClass struct{}

// TaskListItemClass is an extension that renders list items with CSS class
// "task-list-item" if it contains a checkbox
var TaskListItemClass = &taskListItemClass{}

// Extend extends the markdown renderer
func (e *taskListItemClass) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewHTMLRenderer(), 500),
	))
}

// Interface guards
var (
	_ renderer.NodeRenderer = (*HTMLRenderer)(nil)
	_ goldmark.Extender     = (*taskListItemClass)(nil)
)
