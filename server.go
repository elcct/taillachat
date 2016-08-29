package main

import (
	"flag"
	"github.com/elcct/taillachat/api"
	"github.com/elcct/taillachat/system"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
	"net/http"
)

// use is a middleware chainer
func use(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func main() {
	flag.Parse()
	defer glog.Flush()

	err := system.Init()
	if err != nil {
		glog.Fatal(err)
	}

	system.LoadTemplates()

	api.Template = system.CurrentApplication.Template
	api.MediaContent = system.CurrentApplication.Configuration.PublicPath + "/uploads/"

	router := mux.NewRouter()
	// Setup static files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(system.CurrentApplication.Configuration.PublicPath))))

	router.Path("/chat/{any:.*}").Handler(sockjs.NewHandler("/chat", sockjs.DefaultOptions, api.Chat))

	router.Handle("/", use(http.HandlerFunc(api.Index), system.Templates))
	router.Handle("/terms", use(http.HandlerFunc(api.Terms), system.Templates))
	router.Handle("/privacy", use(http.HandlerFunc(api.Privacy), system.Templates))

	glog.Error(http.ListenAndServe(system.CurrentApplication.Configuration.Bind, router))
}
