package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<title>{{ .Title }}</title>
</head>
<body>
{{ .Body }}
</body>
</html>
`
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

var (
	tmpFile string
	server  = NewServer()
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")

	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error: ", err)
		}
	}()

	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		os.Exit(1)
	}
	go watchFile(watcher, *filename, *tFname)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c

	if err := server.s.Shutdown(nil); err != nil {
		return
	}
	if err := watcher.Close(); err != nil {
		fmt.Fprintln(os.Stderr, "Error closing watcher: ", err)
	}
	if tmpFile != "" {
		removeTmpFile(tmpFile)
	}
}

func run(filename, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	if tmpFile != "" {
		removeTmpFile(tmpFile)
	}

	tmp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	defer tmp.Close()
	outName := tmp.Name()
	tmpFile = outName
	fmt.Fprintln(out, outName)

	err = saveHTML(outName, htmlData)
	if err != nil {
		return err
	}
	if skipPreview {
		return nil
	}

	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	var buf bytes.Buffer
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	if err := t.Execute(&buf, c); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func saveHTML(outName string, data []byte) error {
	return os.WriteFile(outName, data, 0644)
}

func preview(fname string) error {
	err := server.reload(fname)
	if err != nil {
		return err
	}

	return nil
}

func removeTmpFile(file string) error {
	return os.Remove(file)
}
