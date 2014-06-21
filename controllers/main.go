package controllers

import (			
	//"github.com/golang/glog"
	"net/http"

	"github.com/zenazn/goji/web"
	"html/template"
	"github.com/elcct/taillachat/helpers"
	"github.com/elcct/taillachat/system"
)

type MainController struct {
	system.Controller
}

func (controller *MainController) Index(c web.C, r *http.Request) (string, int) {	
	t := controller.GetTemplate(c)


	widgets := helpers.Parse(t, "home", nil)

	c.Env["Content"] = template.HTML(widgets)

	c.Env["Title"] = "Tailla Chat - Best UK Chat"
	c.Env["SocketURL"] = "/chat"
	
	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

func (controller *MainController) Terms(c web.C, r *http.Request) (string, int) {	
	t := controller.GetTemplate(c)

	c.Env["Title"] = "Tailla Chat - Best UK Chat - Terms"

	widgets := helpers.Parse(t, "terms", nil)

	c.Env["Content"] = template.HTML(widgets)
	
	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

func (controller *MainController) Privacy(c web.C, r *http.Request) (string, int) {	
	t := controller.GetTemplate(c)

	c.Env["Title"] = "Tailla Chat - Best UK Chat - Privacy"

	widgets := helpers.Parse(t, "privacy", nil)

	c.Env["Content"] = template.HTML(widgets)
	
	return helpers.Parse(t, "main", c.Env), http.StatusOK
}

