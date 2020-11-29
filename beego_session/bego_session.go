package beego_session

import (
	"context"
	"encoding/json"
	"fmt"
	beego_session "github.com/astaxie/beego/session"
	"github.com/code-and-chill/session"
	"net/http"
	"net/http/httptest"
)

func New(engine *beego_session.Manager) *BeegoSession {
	return &BeegoSession{
		Manager: engine,
	}
}

type BeegoSession struct {
	*beego_session.Manager
}

func (beegosession BeegoSession) getSession(w http.ResponseWriter, req *http.Request) (beego_session.Store, error) {
	return beegosession.Manager.SessionStart(w, req)
}

func (beegosession BeegoSession) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	sess, _ := beegosession.getSession(w, req)
	defer sess.SessionRelease(w)

	if str, ok := value.(string); ok {
		return sess.Set(key, str)
	}
	result, _ := json.Marshal(value)
	return sess.Set(key, string(result))
}

func (beegosession BeegoSession) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	sess, _ := beegosession.getSession(w, req)
	defer sess.SessionRelease(w)

	result := sess.Delete(key)
	if result != nil {
		return fmt.Sprint(result)
	}
	return ""
}

// Get value from session data
func (beegosession BeegoSession) Get(req *http.Request, key string) string {
	sess, _ := beegosession.getSession(httptest.NewRecorder(), req)

	result := sess.Get(key)
	if result != nil {
		return fmt.Sprint(result)
	}
	return ""
}

// Flash add flash message to session data
func (beegosession BeegoSession) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := beegosession.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return beegosession.Add(w, req, "_flashes", messages)
}

func (beegosession BeegoSession) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	beegosession.PopLoad(w, req, "_flashes", &messages)
	return messages
}

func (beegosession BeegoSession) Load(req *http.Request, key string, result interface{}) error {
	value := beegosession.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

func (beegosession BeegoSession) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := beegosession.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

func (beegosession BeegoSession) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), "gorilla_writer", w)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
