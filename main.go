package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	header = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>Markdown Preview Tool</title>
</head>
<body>
`
	footer = `
</body>
</html>
`
)

func main() {
	filename := flag.String("filename", "", "Markdown file to preview")

	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*filename, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, out io.Writer) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	tmp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
  defer tmp.Close()
	outName := tmp.Name()
	fmt.Fprintln(out, outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	var buf bytes.Buffer
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	buf.WriteString(header)
	buf.Write(body)
	buf.WriteString(footer)

	return buf.Bytes()
}

func saveHTML(outName string, data []byte) error {
	return os.WriteFile(outName, data, 0644)
}
