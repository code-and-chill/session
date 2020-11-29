package manager

import (
	"github.com/code-and-chill/middlewares"
	"github.com/code-and-chill/session"
	"github.com/code-and-chill/session/gorilla"
	"net/http"

	"github.com/gorilla/sessions"
)

var SessionManager session.IManager = gorilla.New("_session", sessions.NewCookieStore([]byte("secret")))

func init() {
	middlewares.Use(middlewares.Middleware{
		Name: "session",
		Handler: func(handler http.Handler) http.Handler {
			return SessionManager.Middleware(handler)
		},
	})
}
