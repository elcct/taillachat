package system

import (
	"github.com/elcct/taillachat/models"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"net/http"
)

// ApplyTemplates makes sure templates are stored in the context
func (application *Application) ApplyTemplates(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["Template"] = application.Template
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// ApplySessions makes sure controllers can have access to session
func (application *Application) ApplySessions(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := application.Store.Get(r, "session")
		c.Env["Session"] = session
		h.ServeHTTP(w, r)
		context.Clear(r)
	}
	return http.HandlerFunc(fn)
}

// ApplyAuth makes sure user object is in the context
func (application *Application) ApplyAuth(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		session := c.Env["Session"].(*sessions.Session)
		if userID, ok := session.Values["User"].(uint); ok {

			user := &models.User{
				ID: userID,
			}

			c.Env["User"] = user
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
