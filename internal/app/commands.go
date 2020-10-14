package app

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"sidus.io/md3pdf/internal/generated/assets"
	"sidus.io/md3pdf/internal/pkg/latex"
	"strings"
)

const (
	clsName = "md3pdf"
	clsFileName = clsName + ".cls"
	copyCommandName = "cp"
)

var Md3PdfCommand = &cobra.Command{
	Use:   "md3pdf",
	Short: "Convert Markdown to PDF via LaTeX",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	input := args[0]
	inputFile, err := ioutil.ReadFile(input)
	if err != nil {
		panic(err)
	}

	l := latex.NewRenderer()
	md := goldmark.New(
		goldmark.WithExtensions(extension.Table),
		goldmark.WithRenderer(l),
	)
	var buf bytes.Buffer
	err = md.Convert(inputFile, &buf)
	if err != nil {
		panic(err)
	}

	fmt.Print(buf.String())

	err = texToPdf(buf.Bytes(), strings.Split(input, ".")[0], l.Figures)
	if err != nil {
		fmt.Printf("Couldn't generate pdf from tex source")
		panic(err)
	}

}

func texToPdf(texBytes []byte, fileName string, figures []string) error {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "md3pdf-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	texFileName := fmt.Sprintf("%s/%s.tex", tmpDir, fileName)
	err = ioutil.WriteFile(texFileName, texBytes, 0644)
	if err != nil {
		return err
	}

	asset, err := assets.Asset(clsFileName)
	if err != nil {
		return err
	}
	clsFileName := fmt.Sprintf("%s/%s", tmpDir, clsFileName)
	err = ioutil.WriteFile(clsFileName, asset, 0644)
	if err != nil {
		return err
	}

	err = moveLocalFigures(tmpDir, figures)
	if err != nil {
		return err
	}

	err = getRemoteFigures(tmpDir, figures)
	if err != nil {
		return err
	}

	pdflatexCmd := exec.Command("pdflatex", texFileName)
	pdflatexCmd.Dir = tmpDir
	err = pdflatexCmd.Run()
	if err != nil {
		return err
	}
	err = exec.Command(copyCommandName, fmt.Sprintf("%s/%s.pdf", tmpDir, fileName), ".").Run()
	if err != nil {
		return err
	}
	return nil
}

func moveLocalFigures(dir string, figures []string) error {
	for _, fig := range figures {
		dest, err := url.Parse(fig)
		if err != nil || dest.Host == "" || dest.Scheme == "" {
			err := exec.Command(copyCommandName, fig, dir).Run()
			if err != nil {
				return errors.Wrapf(err, "couldn't copy file %s", fig)
			}
		}

	}
	return nil
}

func getRemoteFigures(dir string, figures []string) error {
	for _, fig := range figures {
		dest, err := url.Parse(fig)
		if err != nil || dest.Host == "" || dest.Scheme == "" {
			continue
		}
		wgetCommand := exec.Command("wget", fig)
		wgetCommand.Dir = dir
		err = wgetCommand.Run()
		if err != nil {
			return errors.Wrapf(err, "couldn't download %s", fig)
		}
	}
	return nil
}