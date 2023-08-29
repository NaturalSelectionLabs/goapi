package main

import (
	"bytes"
	"html/template"
	"os"

	"github.com/NaturalSelectionLabs/goapi/lib/openapi"
)

// create a golang file "response_status_code.go" in the current directory.
func main() {
	out := "package goapi\n"

	for _, code := range openapi.StatusCode(0).Values() {
		out += format(`
type Status{{.}} struct{}

func (Status{{.}}) statusCode() int {
	return {{.}}
}
`, code)
	}

	f, err := os.OpenFile("response-status-code.go", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644) //nolint: gofumpt
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(out)
	if err != nil {
		panic(err)
	}
}

// render go template.
func format(tpl string, value string) string {
	var buf bytes.Buffer

	t := template.Must(template.New("").Parse(tpl))

	err := t.Execute(&buf, value)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
