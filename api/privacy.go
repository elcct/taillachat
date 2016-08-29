package api

import (
	"github.com/elcct/taillachat/helpers"
	"html/template"
	"net/http"
)

// Privacy serves private page
func Privacy(w http.ResponseWriter, r *http.Request) {
	t := r.Context().Value("template").(*template.Template)
	widgets := helpers.Parse(t, "privacy", nil)

	data := map[string]interface{}{}
	data["Content"] = template.HTML(widgets)

	data["Title"] = "Tailla Chat - Best UK Chat - Privacy"

	w.Write([]byte(helpers.Parse(t, "main", data)))
}
