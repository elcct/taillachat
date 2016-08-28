package system

import (
	"github.com/gorilla/context"
	"net/http"
)

// Templates adds templates to the context
func Templates(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context.Set(r, "template", CurrentApplication.Template)
		inner.ServeHTTP(w, r)
	})
}
