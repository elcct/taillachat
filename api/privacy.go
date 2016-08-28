package api

import (
	"github.com/elcct/taillachat/helpers"
	"github.com/gorilla/context"
	"html/template"
	"net/http"
)

// Privacy serves private page
func Privacy(w http.ResponseWriter, r *http.Request) {
	t := context.Get(r, "template").(*template.Template)
	widgets := helpers.Parse(t, "privacy", nil)

	data := map[string]interface{}{}
	data["Content"] = template.HTML(widgets)

	data["Title"] = "Tailla Chat - Best UK Chat - Privacy"

	w.Write([]byte(helpers.Parse(t, "main", data)))
}
