package latex

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark/ast"
)

func (l *latex) renderInline(node ast.Node) error {
	var err error
	switch node.Kind() {
	case ast.KindText:
		err = l.write(string(node.Text(l.source)))
	case ast.KindEmphasis:
		err = l.renderEmphasis(node.(*ast.Emphasis))
	default:
		return fmt.Errorf("redering not implemented: inline kind: %s", node.Kind().String())
	}
	if err != nil {
		return errors.Wrapf(err, "coudln't render inline kind: %s", node.Kind().String())
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