package session

import (
	"html/template"
	"net/http"
)

type IManager interface {
	Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error
	Get(req *http.Request, key string) string
	Pop(w http.ResponseWriter, req *http.Request, key string) string

	Flash(w http.ResponseWriter, req *http.Request, message Message) error
	Flashes(w http.ResponseWriter, req *http.Request) []Message

	Load(req *http.Request, key string, result interface{}) error
	PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error

	Middleware(handler http.Handler) http.Handler
}

type Message struct {
	Message template.HTML
	Type    string
}
