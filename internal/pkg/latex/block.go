package latex

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/pkg/errors"
)

func (l *latex) renderBlock(node ast.Node) error {
	var err error
	switch node.Kind() {
	case ast.KindHeading:
		err = l.renderHeading(node.(*ast.Heading))
	case ast.KindParagraph:
		err = l.renderParagraph(node.(*ast.Paragraph))
	case ast.KindList:
		err = l.renderList(node.(*ast.List))
	case ast.KindListItem:
		err = l.renderListItem(node.(*ast.ListItem))
	case ast.KindTextBlock:
		err = l.render(node.FirstChild())
	case ast.KindFencedCodeBlock:
		err = l.renderFencedCodeBlock(node.(*ast.FencedCodeBlock))
	default:
		return fmt.Errorf("not implemented block kind %s", node.Kind().String())
	}

	if err != nil {
		return errors.Wrapf(err, "couldn't render block kind: %s", node.Kind().String())
	}

	return nil
}


func (l *latex) renderHeading(h *ast.Heading) error {
	var err error
	switch h.Level {
	case 1:
		err = l.write("\\chapter{")
	case 2:
		err = l.write("\\section{")
	case 3:
		err = l.write("\\subsection{")
	case 4:
		err = l.write("\\subsubsection{")
	case 5:
		err = l.write("\\paragraph{")
	case 6:
		err = l.write("\\subparagraph{")
	default:
		return fmt.Errorf("not implemented heading level %d", h.Level)
	}
	if err != nil {
		return errors.Wrapf(err, "couldn't render heading level: %d", h.Level)
	}

	err = l.render(h.FirstChild())
	if err != nil {
		return errors.Wrapf(err, "couldn't render heading level: %d", h.Level)
	}
	err = l.writeln("}")
	if err != nil {
		return errors.Wrapf(err, "couldn't render heading level: %d", h.Level)
	}
	return nil
}

func (l *latex) renderParagraph(p *ast.Paragraph) error {
	err := l.render(p.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render paragraph")
	}
	err = l.write("\n\n")
	if err != nil {
		return errors.Wrap(err, "couldn't render paragraph")
	}
	return nil
}

func (l *latex) renderList(list *ast.List) error {
	var listType string
	if list.IsOrdered() {
		listType = "enumerate"
	} else {
		listType = "itemize"
	}
	err := l.writefln("\\begin{%s}", listType)
	if err != nil {
		return errors.Wrap(err, "couldn't render list")
	}

	l.indent++
	err = l.render(list.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render list")
	}
	l.indent--

	err = l.writefln("\\end{%s}", listType)
	return nil
}

func (l *latex) renderListItem(item *ast.ListItem) error {
	err := l.write("\\item ")
	if err != nil {
		return errors.Wrap(err, "couldn't render listitem")
	}
	err = l.render(item.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render listitem")
	}
	err = l.write("\n")
	if err != nil {
		return errors.Wrap(err, "couldn't render listitem")
	}
	return nil
}

func (l *latex) renderFencedCodeBlock(block *ast.FencedCodeBlock) error {
	lang := block.Language(l.source)
	err := l.write("\\begin{lstlisting}")
	if err != nil {
		return errors.Wrap(err, "coudln't render FencedCodeBlock")
	}
	err = l.writefln( "[language=%s]", lang)
	if err != nil {
		return errors.Wrap(err, "coudln't render FencedCodeBlock")
	}
	l.writeLines(block)
	if err != nil {
		return errors.Wrap(err, "coudln't render FencedCodeBlock")
	}
	err = l.writeln("\\end{lstlisting}")
	return nil
}