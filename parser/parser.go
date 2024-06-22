package parser

import (
	"bytes"
	"html/template"
)

func ParseTemplateString(fileName, tmpl string, data interface{}) (string, error) {
	t, err := template.New(fileName).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
