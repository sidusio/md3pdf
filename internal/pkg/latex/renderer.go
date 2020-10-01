package latex

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"io"
)

type latex struct{

}

func NewRenderer() *latex {
	return &latex{}
}

func (l *latex) Render(w io.Writer, source []byte, node ast.Node) error {
	for n := node; n != nil; n = n.NextSibling() {
		switch n.Type() {
		case ast.TypeDocument:
			err := l.renderDocument(w, source, n)
			if err != nil {
				return err
			}
		case ast.TypeBlock:
			err := l.renderBlock(w, source, n)
			if err != nil {
				return err
			}
		case ast.TypeInline:
			err := l.renderInline(w, source, n)
			if err != nil {
				return nil
			}
		default:
			return fmt.Errorf("not implemented %d", n.Type())
		}
	}
	return nil
}

func (l *latex) AddOptions(option ...renderer.Option) {
	panic("implement me")
}

func (l *latex) renderDocument(w io.Writer, source []byte, node ast.Node) error {
	_, _ = fmt.Fprintf(w, "\\documentclass{md3pdf}\n")
	_, _ = fmt.Fprintf(w, "\\begin{document}\n")

	err := l.Render(w, source, node.FirstChild())
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(w, "\\end{document}\n")
	return nil
}

func (l *latex) renderBlock(w io.Writer, source []byte, node ast.Node) error {
	switch node.Kind() {
	case ast.KindHeading:
		err := l.renderHeading(w, source, node.(*ast.Heading))
		if err != nil {
			return err
		}
	case ast.KindParagraph:
		err := l.renderParagraph(w, source, node.(*ast.Paragraph))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("not implemented block kind %s", node.Kind().String())
	}
	return nil
}

func (l *latex) renderInline(w io.Writer, source []byte, node ast.Node) error {
	switch node.Kind() {
	case ast.KindText:
		_,_ = fmt.Fprint(w, string(node.Text(source)))
	}
	return nil
}

func (l *latex) renderHeading(w io.Writer, source []byte, h *ast.Heading) error {
	switch h.Level {
	case 1:
		_, _ = fmt.Fprint(w, "\\section{")
	case 2:
		_, _ = fmt.Fprint(w, "\\subsection{")
	default:
		return fmt.Errorf("not implemented heading level %d", h.Level)
	}
	err := l.Render(w, source, h.FirstChild())
	if err != nil {
		return err
	}
	_, _ = fmt.Fprint(w, "}\n")
	return nil
}

func (l *latex) renderParagraph(w io.Writer, source []byte, p *ast.Paragraph) error {
	err := l.Render(w, source, p.FirstChild())
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(w, "\n\n")
	return nil
}