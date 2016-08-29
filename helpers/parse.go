package helpers

import (
	"bytes"
	"html/template"
)

// Parse parses template of `name` and fills it with `data`
// TODO: we ignore the error on purpose, that should be refactored later on
func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}
