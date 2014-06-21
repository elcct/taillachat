package helpers

import (
	"html/template"
	"bytes"
	"fmt"
)

func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	err := t.ExecuteTemplate(&doc, name, data)
	if err != nil {
		fmt.Println(err)
	}
	return doc.String()
}

