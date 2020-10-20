package app

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sidus.io/md3pdf/internal/generated/assets"
	"sidus.io/md3pdf/internal/pkg/latex"
	"strings"
	"time"
)

const (
	clsName             = "md3pdf"
	clsFileName         = clsName + ".cls"
	pdfLatexCommandName = "pdflatex"
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
		fmt.Printf("Couldn't generate pdf from tex source\n")
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

	err = getFigures(tmpDir, figures)
	if err != nil {
		return err
	}

	err = convertFigures(tmpDir, figures)
	if err != nil {
		return err
	}

	pdflatexCmd := exec.Command(pdfLatexCommandName, texFileName)
	pdflatexCmd.Dir = tmpDir
	err = pdflatexCmd.Run()
	if err != nil {
		return err
	}
	err = copyFile(fmt.Sprintf("%s/%s.pdf", tmpDir, fileName), ".")
	if err != nil {
		return err
	}
	return nil
}

func getFigures(tmpDir string, figures []string) error {
	for _, fig := range figures {
		dest, err := url.Parse(fig)
		if err != nil || dest.Host == "" || dest.Scheme == "" {
			err = getLocalFigure(tmpDir, fig)
		} else {
			err = getRemoteFigure(tmpDir, fig)
		}
		if err != nil {
			return errors.Wrapf(err, "couldn't get figure %s", fig)
		}
	}
	return nil
}

func getLocalFigure(dir string, fig string) error {
	return copyFile(fig, dir)
}

func getRemoteFigure(dir string, fig string) error {
	_, err := url.Parse(fig)
	if err != nil {
		return err
	}
	err = downloadFigure(dir, fig)
	if err != nil {
		return errors.Wrapf(err, "couldn't download %s", fig)
	}
	return nil
}

func convertFigures(dir string, figures []string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(dir)
	if err != nil {
		return errors.Wrapf(err, "couldn't change directory to %s", dir)
	}
	defer must(os.Chdir, wd, "Couldn't change directory")

	for _, fig := range figures {
		parts := strings.Split(fig, string(filepath.Separator))
		fileName := parts[len(parts)-1]
		err := convertFigure(fileName)
		if err != nil {
			return errors.Wrapf(err, "couldn't convert figure %s", fig)
		}
	}
	return nil
}

func convertFigure(fig string) error {
	err := exec.Command("magick", fig, fig+".png").Run()
	if err != nil {
		return errors.Wrapf(err, "magick couldn't convert %s to %s")
	}
	return nil
}

func downloadFigure(dir string, fig string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	defer must(os.Chdir, wd, "Couldn't change directory")

	// Build fileName from fullPath
	fileURL, err := url.Parse(fig)
	if err != nil {
		return err
	}
	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
		Timeout: time.Second * 10,
	}

	resp, err := client.Get(fig)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	defer file.Close()
	if err != nil {
		return err
	}

	return nil
}

func copyFile(source string, destination string) error {
	pathStack := strings.Split(source, string(os.PathSeparator))
	fileName := pathStack[len(pathStack)-1]

	src, err := os.Open(source)
	if err != nil {
		return err
	}

	dest, err := os.Open(destination)
	if err != nil {
		return err
	}
	fi, err := dest.Stat()
	if err != nil {
		return err
	}

	if fi.IsDir() {
		dest, err = os.Create(destination + string(os.PathSeparator) + fileName)
		if err != nil {
			return err
		}
	}

	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}

	return nil
}

func must(f func(string) error, s string, msg string) {
	err := f(s)
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}