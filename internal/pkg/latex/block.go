package latex

import (
	"fmt"
	"github.com/yuin/goldmark/ast"
	"github.com/pkg/errors"
	ast2 "github.com/yuin/goldmark/extension/ast"
	"strings"
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
	case ast2.KindTable:
		err = l.renderTable(node.(*ast2.Table))
	case ast2.KindTableHeader:
		err = l.renderTableHeader(node.(*ast2.TableHeader))
	case ast2.KindTableRow:
		err = l.renderTableRow(node.(*ast2.TableRow))
	case ast2.KindTableCell:
		err = l.renderTableCell(node.(*ast2.TableCell))
	case ast.KindBlockquote:
		err = l.renderBlockQuote(node.(*ast.Blockquote))
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
	if string(lang) != "" {
		err = l.writefln( "[language=%s]", lang)
		if err != nil {
			return errors.Wrap(err, "coudln't render FencedCodeBlock")
		}
	}

	l.writeLines(block)
	if err != nil {
		return errors.Wrap(err, "coudln't render FencedCodeBlock")
	}
	err = l.writeln("\\end{lstlisting}")
	return nil
}

func (l *latex) renderTable(table *ast2.Table) error {
	err := l.writeln("\\begin{table}[H]")
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}
	l.indent++

	err = l.writeln("\\rowcolors{2}{white!80!black!50}{white!70!black!40}")
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}

	alignments := make([]string, len(table.Alignments))
	for i, alignment := range table.Alignments {
		switch alignment {
		case ast2.AlignCenter:
			alignments[i] = "c"
		case ast2.AlignRight:
			alignments[i] = "r"
		default:
			alignments[i] = "l"
		}
	}

	err = l.writefln("\\begin{tabular}{|%s|}", strings.Join(alignments, "|"))
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}
	l.indent++

	err = l.writeln("\\hline")
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}

	err = l.render(table.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}

	l.indent--
	err = l.writeln("\\end{tabular}")
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}

	l.indent--
	err = l.writeln("\\end{table}")
	if err != nil {
		return errors.Wrap(err, "couldn't render table")
	}
	return nil
}

func (l *latex) renderTableHeader(header *ast2.TableHeader) error {
	l.isTableHeader = true
	err := l.renderTableRow(header)
	if err != nil {
		return errors.Wrap(err, "couldn't render table header")
	}
	l.isTableHeader = false
	return nil
}

func (l *latex) renderTableRow(row ast.Node) error {
	l.firstTableRowCellWritten = false // Set to true in renderTableCell

	err := l.render(row.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render table row")
	}

	err = l.writeln("\\\\ \\hline")
	if err != nil {
		return errors.Wrap(err, "couldn't render table row")
	}
	return nil
}

func (l *latex) renderTableCell(row *ast2.TableCell) error {
	var err error
	if l.firstTableRowCellWritten {
		err = l.write("& ")
		if err != nil {
			return errors.Wrap(err, "couldn't render table cell")
		}
	}
	l.firstTableRowCellWritten = true // Set to false in renderTableRow

	err = l.render(row.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render table cell")
	}
	return nil
}

func (l *latex) renderBlockQuote(blockquote *ast.Blockquote) error {
	_ = l.write("\n")
	_ = l.writeln("\\begin{myquote}")

	err := l.render(blockquote.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render blockquote")
	}

	_ = l.writef("\\end{myquote}\n\n")
	return nil
}
