package latex

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"io"
	"strings"
)

type latex struct{
	indent int
	w io.Writer
	source []byte
}

func NewRenderer() *latex {
	return &latex{}
}

func (l *latex) Render(w io.Writer, source []byte, node ast.Node) error {
	node.Dump(source, 1)
	l.source = source
	l.w = w
	return l.render(node)
}

func (l *latex) render(node ast.Node) error {

	for n := node; n != nil; n = n.NextSibling() {
		var err error
		switch n.Type() {
		case ast.TypeDocument:
			err = l.renderDocument(n)
		case ast.TypeBlock:
			err = l.renderBlock(n)
		case ast.TypeInline:
			err = l.renderInline(n)
		default:
			return fmt.Errorf("not implemented %d", n.Type())
		}
		if err != nil {
			return errors.Wrapf(err, "coudln't render type: %d", n.Type())
		}
	}
	return nil
}



func (l *latex) AddOptions(option ...renderer.Option) {
	panic("implement me")
}

func (l *latex) renderDocument( node ast.Node) error {
	_ = l.writeln("\\documentclass{md3pdf}")
	_ = l.writeln("\\begin{document}")

	err := l.render(node.FirstChild())
	if err != nil {
		return err
	}

	_ = l.writeln("\\end{document}")
	return nil
}

func (l *latex) write(s string) error {
	tabs := strings.Repeat("\t", l.indent)
	_, err := fmt.Fprintf(l.w, "%s%s", tabs, s)
	return err
}

func (l *latex) writef(format string, s ...interface{}) error {
	return l.write(fmt.Sprintf(format, s...))
}

func (l *latex) writeln(s string) error {
	return l.write(fmt.Sprintf("%s\n", s))
}

func (l *latex) writefln(format string, s ...interface{}) error {
	return l.writef(format + "\n", s...)
}

func (l *latex) writeLines(n ast.Node) {
	line := n.Lines().Len()
	for i := 0; i < line; i++ {
		line := n.Lines().At(i)
		_, _ = fmt.Fprint(l.w, string(line.Value(l.source)))
	}
}