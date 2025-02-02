package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

type res map[string]error

type Result struct {
	w       io.Writer
	Out     string
	Success int
	Abort   int
	URLS    res
}

// HandleEvent implements pubsub.Subscriber.
func (r *Result) HandleEvent(event Event) {
	switch e := event.(type) {
	case EventStart:
	case EventProgress:
	case EventEnd:
		r.Success++
	case EventRetry:
	case EventAbort:
		r.URLS[e.URL] = e.Err
		r.Abort++
	default:
		panic(fmt.Sprintf("unexpected main.Event: %#v", event))
	}
}

func NewResult(w io.Writer, outDir string) *Result {
	return &Result{
		w:       w,
		Out:     outDir,
		URLS:    make(res),
		Success: 0,
		Abort:   0,
	}
}

const format = `
Stored {{.Success}} files to {{.Out}}.
{{ if .Abort }}Aborted {{ .Abort }} urls:
Error: {{ range $key, $err := .URLS }} 
	- {{$key}}: {{ PrettyError $err }}
{{ end }}{{- end}}`

// Stored $n files to $out.
// Aborted $url:
//
// Error:
//   - url1: $error1
//   - url2: $error2
func (r *Result) Show() {
	tmpl, err := template.New("test").
		Funcs(template.FuncMap{"PrettyError": PrettyError}).
		Parse(format)
	if err != nil {
		panic(err)
	}

	tmpl.Execute(os.Stdout, r)
}

func PrettyError(e error) string {
	str := e.Error()
	deduped := make(map[string]bool)

	for _, line := range strings.Split(str, "\n") {
		deduped[line] = true
	}

	ret := ""
	for key := range deduped {
		ret += key
	}
	return ret
}
