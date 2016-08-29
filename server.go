package main

import (
	"flag"
	"github.com/elcct/taillachat/api"
	"github.com/elcct/taillachat/system"
	"github.com/golang/glog"
	"github.com/gorilla/context"
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

	api.Template = system.CurrentApplication.Template
	api.MediaContent = system.CurrentApplication.Configuration.PublicPath + "/uploads/"

	router := mux.NewRouter()
	// Setup static files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(system.CurrentApplication.Configuration.PublicPath))))

	router.Path("/chat/{any:.*}").Handler(sockjs.NewHandler("/chat", sockjs.DefaultOptions, api.Chat))

	router.Handle("/", use(http.HandlerFunc(api.Index), templates))
	router.Handle("/terms", use(http.HandlerFunc(api.Terms), templates))
	router.Handle("/privacy", use(http.HandlerFunc(api.Privacy), templates))

	glog.Error(http.ListenAndServe(system.CurrentApplication.Configuration.Bind, router))
}

// Templates adds templates to the context
func templates(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "template", system.CurrentApplication.Template)
		inner.ServeHTTP(w, r)
	})
}
