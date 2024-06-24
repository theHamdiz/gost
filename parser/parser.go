package parser

import (
	"bytes"
	ht "html/template"
	"strings"
	tt "text/template"
)

func ParseTemplateStringAsText(fileName, tmpl string, data interface{}) (string, error) {
	// Determine whether to use html/template or text/template
	if strings.HasSuffix(fileName, ".html") || strings.HasSuffix(fileName, ".htm") || strings.HasSuffix(fileName, ".tmpl") || strings.HasSuffix(fileName, ".gohtml") {
		htmpl, err := ht.New(fileName).Parse(tmpl)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err := htmpl.Execute(&buf, data); err != nil {
			return "", err
		}

		return buf.String(), nil
	} else {
		ttmpl, err := tt.New(fileName).Parse(tmpl)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err := ttmpl.Execute(&buf, data); err != nil {
			return "", err
		}

		return buf.String(), nil
	}
}

func ParseTemplateStringAsHtml(fileName, tmpl string, data interface{}) (string, error) {
	t, err := ht.New(fileName).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
