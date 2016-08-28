package helpers

import (
	"bytes"
	"fmt"
	"html/template"
)

func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	err := t.ExecuteTemplate(&doc, name, data)
	if err != nil {
		fmt.Println(err)
	}
	return doc.String()
}
