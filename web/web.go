package web

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/YouthBuild-USA/godata/log"
	"github.com/goods/httpbuf"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
)

var store sessions.Store

type HandlerAdder func(string, Handle) *mux.Route
type Handle func(http.ResponseWriter, *http.Request) error

var requestLog = log.New("request")

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

type RequestInitializer func(*http.Request) error

var initializers = make([]RequestInitializer, 0, 0)

func AddRequestInitializer(initializer RequestInitializer) {
	initializers = append(initializers, initializer)
}

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestLog.Info("%v %v", r.Method, r.RequestURI)

	for _, initializer := range initializers {
		initializer(r)
	}

	buffer := new(httpbuf.Buffer)

	err := h(buffer, r)
	if err != nil {
		fmt.Println(err)
		if webErr, ok := err.(WebError); ok {
			http.Error(w, webErr.Message, webErr.Code)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	err = Session(r).Save(r, buffer)
	context.Clear(r)

	buffer.Apply(w)
}

type WebError struct {
	Err     error
	Message string
	Code    int
}

func (err WebError) Error() string {
	return fmt.Sprintf("%v: %v", err.Code, err.Message)

}

// CreateSessionStore initializes the user session store for the system
func CreateSessionStore(key string) {
	store = sessions.NewCookieStore([]byte(key))
}

// Session returns the user session
func Session(r *http.Request) *sessions.Session {
	session, _ := store.Get(r, "session-name")
	return session
}

const csrfLength = 20

func GenerateCSRF(r *http.Request) string {
	b := make([]byte, csrfLength)
	io.ReadFull(rand.Reader, b)
	token := hex.EncodeToString(b)
	session := Session(r)
	session.Values["csrf"] = token
	return token
}

func CSRF(r *http.Request) string {
	session := Session(r)
	if token, ok := session.Values["csrf"].(string); ok {
		return token
	}
	return GenerateCSRF(r)
}

func ValidateCSRF(r *http.Request, token string) bool {
	session := Session(r)
	if last, ok := session.Values["csrf"].(string); ok {
		return last == token
	}
	return false
}
