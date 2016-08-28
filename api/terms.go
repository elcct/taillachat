package api

import (
	"github.com/elcct/taillachat/helpers"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)

// Terms serves terms page
func Terms(w http.ResponseWriter, r *http.Request) {
	t := context.Get(r, "template").(*template.Template)
	widgets := helpers.Parse(t, "terms", nil)

	data := map[string]interface{}{}
	data["Content"] = template.HTML(widgets)

	data["Title"] = "Tailla Chat - Best UK Chat - Terms"

	w.Write([]byte(helpers.Parse(t, "main", data)))
}
