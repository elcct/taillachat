package api

import (
	"github.com/elcct/taillachat/helpers"
	"html/template"
	"net/http"
)

// Terms serves terms page
func Terms(w http.ResponseWriter, r *http.Request) {
	t := r.Context().Value("template").(*template.Template)
	widgets := helpers.Parse(t, "terms", nil)

	data := map[string]interface{}{}
	data["Content"] = template.HTML(widgets)

	data["Title"] = "Tailla Chat - Best UK Chat - Terms"

	w.Write([]byte(helpers.Parse(t, "main", data)))
}
