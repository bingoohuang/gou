package lang

import (
	"bytes"
	"text/template"
)

// TmplRenderText renders the tmpl template with data.
func TmplRenderText(tmpl string, data interface{}) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}

	w := &bytes.Buffer{}

	if err = t.Execute(w, data); err != nil {
		return "", err
	}

	return w.String(), nil
}
