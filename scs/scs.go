package scs

import (
	"encoding/json"
	"github.com/code-and-chill/session"
	"net/http"

	"github.com/alexedwards/scs"
)

func New(manager *scs.Manager) *SCS {
	return &SCS{Manager: manager}
}

type SCS struct {
	*scs.Manager
}

func (scs SCS) Add(w http.ResponseWriter, req *http.Request, key string, value interface{}) error {
	scssession := scs.Manager.Load(req)

	if str, ok := value.(string); ok {
		return scssession.PutString(w, key, str)
	}
	result, _ := json.Marshal(value)
	return scssession.PutString(w, key, string(result))
}

func (scs SCS) Pop(w http.ResponseWriter, req *http.Request, key string) string {
	scssession := scs.Manager.Load(req)
	result, _ := scssession.PopString(w, key)
	return result
}

func (scs SCS) Get(req *http.Request, key string) string {
	scssession := scs.Manager.Load(req)
	result, _ := scssession.GetString(key)
	return result
}

func (scs SCS) Flash(w http.ResponseWriter, req *http.Request, message session.Message) error {
	var messages []session.Message
	if err := scs.Load(req, "_flashes", &messages); err != nil {
		return err
	}
	messages = append(messages, message)
	return scs.Add(w, req, "_flashes", messages)
}

func (scs SCS) Flashes(w http.ResponseWriter, req *http.Request) []session.Message {
	var messages []session.Message
	scs.PopLoad(w, req, "_flashes", &messages)
	return messages
}

func (scs SCS) Load(req *http.Request, key string, result interface{}) error {
	value := scs.Get(req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

func (scs SCS) PopLoad(w http.ResponseWriter, req *http.Request, key string, result interface{}) error {
	value := scs.Pop(w, req, key)
	if value != "" {
		return json.Unmarshal([]byte(value), result)
	}
	return nil
}

func (scs SCS) Middleware(handler http.Handler) http.Handler {
	return scs.Manager.Multi(handler)
}
