package latex

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
	"net/url"
	"strings"
)

func (l *latex) renderInline(node ast.Node) error {
	var err error
	switch node.Kind() {
	case ast.KindText:
		err = l.renderText(node.(*ast.Text))
	case ast.KindEmphasis:
		err = l.renderEmphasis(node.(*ast.Emphasis))
	case ast.KindLink:
		err = l.renderLink(node.(*ast.Link))
	case ast.KindImage:
		err = l.renderImage(node.(*ast.Image))
	default:
		return fmt.Errorf("redering not implemented: inline kind: %s", node.Kind().String())
	}
	if err != nil {
		return errors.Wrapf(err, "coudln't render inline kind: %s", node.Kind().String())
	}
	return nil
}

func (l *latex) renderText(text *ast.Text) error {
	err := l.write(escapeText(string(text.Text(l.source))))
	if err != nil {
		return errors.Wrap(err, "couldn't render text")
	}
	return nil
}

func (l *latex) renderEmphasis(emph *ast.Emphasis) error {
	var levelString string
	switch emph.Level {
	case 1:
		levelString = "emph"
	case 2:
		levelString = "textbf"
	default:
		return fmt.Errorf("not implelemted for level %d", emph.Level)
	}
	err := l.writef("\\%s{", levelString)
	if err != nil {
		return errors.Wrap(err, "couldn't render emphasis")
	}
	err = l.render(emph.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render emphasis")
	}
	err = l.write("}")
	if err != nil {
		return errors.Wrap(err, "couldn't render emphasis")
	}
	return nil
}

func (l *latex) renderLink(link *ast.Link) error {
	err := l.writef("\\href{%s}{", link.Destination)
	if err != nil {
		return errors.Wrap(err, "couldn't render link")
	}
	err = l.render(link.FirstChild())
	if err != nil {
		return errors.Wrap(err, "couldn't render link")
	}
	err = l.write("}")
	if err != nil {
		return errors.Wrap(err, "couldn't render link")
	}
	return nil
}
func (l *latex) renderImage(image *ast.Image) error {
	var file string
	dest, err := url.Parse(string(image.Destination))
	if err != nil || dest.Host == "" || dest.Scheme == "" {
		// Refers to a local file
		file = string(image.Destination)
	} else {
		// refers to an url
		l.Figures = append(l.Figures, dest.String())
		p := strings.Split(dest.Path, "/")
		file = p[len(p)-1]
	}
	if strings.Split(file, ".")[1] == "svg" {
		_ = l.writef("\\includesvg[max width=0.9\\linewidth]{%s}", file)
	} else {
		_ = l.writef("\\includegraphics[max width=0.9\\linewidth]{%s}", file)
	}

	return nil
}