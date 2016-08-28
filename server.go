package main

import (
	"flag"
	"github.com/elcct/taillachat/controllers"
	"github.com/elcct/taillachat/system"
	"github.com/golang/glog"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"net/http"
)

func main() {
	filename := flag.String("config", "config.json", "Path to configuration file")

	flag.Parse()
	defer glog.Flush()

	var application = &system.Application{}

	application.Init(filename)
	application.LoadTemplates()
	application.ConnectToDatabase()

	// Setup static files
	static := web.New()
	static.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(application.Configuration.PublicPath))))

	http.Handle("/assets/", static)

	// That's probably a terrible idea
	controllers.Template = application.Template
	controllers.MediaContent = application.Configuration.PublicPath + "/uploads/"

	http.Handle("/chat/", sockjs.NewHandler("/chat", sockjs.DefaultOptions, controllers.Chat))

	// Apply middleware
	goji.Use(application.ApplyTemplates)
	goji.Use(application.ApplySessions)
	goji.Use(application.ApplyDatabase)
	goji.Use(application.ApplyAuth)

	controller := &controllers.MainController{}

	goji.Get("/", application.Route(controller, "Index"))
	goji.Get("/terms", application.Route(controller, "Terms"))
	goji.Get("/privacy", application.Route(controller, "Privacy"))

	graceful.PostHook(func() {
		application.Close()
	})
	goji.Serve()
}
