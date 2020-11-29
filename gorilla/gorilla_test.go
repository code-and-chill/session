package gorilla_test

import (
	"github.com/code-and-chill/session/gorilla"
	"github.com/code-and-chill/session/test"
	"testing"

	"github.com/gorilla/sessions"
)

func TestAll(t *testing.T) {
	engine := sessions.NewCookieStore([]byte("something-very-secret"))
	manager := gorilla.New("_session", engine)
	test.TestAll(manager, t)
}
