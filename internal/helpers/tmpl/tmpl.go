package tmpl

import (
	"bytes"
	"text/template"
)

func GenerateString(tmpl *template.Template, data any) string {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func GenerateBytes(tmpl *template.Template, data any) []byte {
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
