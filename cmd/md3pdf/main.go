package main

import (
	"os"
	"sidus.io/md3pdf/internal/app"
)

func main() {
	if err := app.Md3PdfCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
