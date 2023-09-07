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
	outTest := `package goapi

import "testing"

func TestStatusCode(t *testing.T) {`

	for _, name := range openapi.StatusCode(0).Values() {
		code, _ := openapi.StatusCodeString(name)
		out += format(`
// Status{{.Name}} {{.Code}}.
type Status{{.Name}} struct{}

func (Status{{.Name}}) statusCode() int {
	return {{.Code}}
}
`, struct {
			Name string
			Code int
		}{name, int(code)})

		outTest += format(`
	Status{{.Name}}{}.statusCode()`, struct {
			Name string
			Code int
		}{name, int(code)})
	}

	outTest += `
}
`

	writeFile("status_code.go", out)
	writeFile("status_code_test.go", outTest)
}

// render go template.
func format(tpl string, value any) string {
	var buf bytes.Buffer

	t := template.Must(template.New("").Parse(tpl))

	err := t.Execute(&buf, value)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func writeFile(name, data string) {
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(data)
	if err != nil {
		panic(err)
	}
}
