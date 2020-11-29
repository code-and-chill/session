package gorilla

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/code-and-chill/session"
	gorillaContext "github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"net/http"
)

func New(sessionName string, store sessions.Store) *Gorilla {
	return &Gorilla{
		SessionName: sessionName,
		Store:       store,
	}
}

type Gorilla struct {
	SessionName string
	Store       sessions.Store
}

func (gorilla Gorilla) getSession(req *http.Request) (*sessions.Session, error) {
	if r, ok := req.Context().Value("gorilla_reader").(*http.Request); ok {
		return gorilla.Store.Get(r, gorilla.SessionName)
	}
	return gorilla.Store.Get(req, gorilla.SessionName)
}

func (gorilla Gorilla) saveSession(w http.ResponseWriter, req *http.Request) {
	if session, err := gorilla.getSession(req); err == nil {
		if err := session.Save(req, w); err != nil {
			fmt.Printf("No error should happen when saving session data, but got %v", err)
		}
	}
}

func (gorilla Gorilla) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	defer gorilla.saveSession(w, req)

	sess, err := gorilla.getSession(req)
	if err != nil {
		return err
	}

	if str, ok := value.(string); ok {
		sess.Values[key] = str
	} else {
		result, _ := json.Marshal(value)
		sess.Values[key] = string(result)
	}

	return nil
}

// Pop value from session data
func (gorilla Gorilla) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	defer gorilla.saveSession(w, req)

	if sess, err := gorilla.getSession(req); err == nil {
		if value, ok := sess.Values[key]; ok {
			delete(sess.Values, key)
			return fmt.Sprint(value)
		}
	}
	return ""
}

// Get value from session data
func (gorilla Gorilla) Get(req *http.Request, key string) string {
	if sess, err := gorilla.getSession(req); err == nil {
		if value, ok := sess.Values[key]; ok {
			return fmt.Sprint(value)
		}
	}
	return ""
}

// Flash add flash message to session data
func (gorilla Gorilla) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := gorilla.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return gorilla.Add(w, req, "_flashes", messages)
}

func (gorilla Gorilla) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	gorilla.PopLoad(w, req, "_flashes", &messages)
	return messages
}

func (gorilla Gorilla) Load(req *http.Request, key string, result interface{}) error {
	value := gorilla.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

func (gorilla Gorilla) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := gorilla.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

func (gorilla Gorilla) Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer gorillaContext.Clear(req)
		ctx := context.WithValue(req.Context(), "gorilla_reader", req)
		handler.ServeHTTP(w, req.WithContext(ctx))
	})
}
