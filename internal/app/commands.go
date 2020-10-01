package app

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"io/ioutil"
	"os"
	"os/exec"
	"sidus.io/md3pdf/internal/pkg/latex"
	"strings"
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

	md := goldmark.New(
		goldmark.WithRenderer(latex.NewRenderer()),
	)
	var buf bytes.Buffer
	err = md.Convert(inputFile, &buf)
	if err != nil {
		panic(err)
	}

	err = texToPdf(buf.Bytes(), strings.Split(input, ".")[0])
	if err != nil {
		fmt.Printf("Couldn't generate pdf from tex source")
		panic(err)
	}


}

func texToPdf(texBytes []byte, fileName string) error {
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

	pdflatexCmd := exec.Command("pdflatex", texFileName)
	pdflatexCmd.Dir = tmpDir
	err = pdflatexCmd.Run()
	if err != nil {
		return err
	}
	err = exec.Command("cp", fmt.Sprintf("%s/%s.pdf", tmpDir, fileName), ".").Run()
	if err != nil {
		return err
	}
	return nil
}