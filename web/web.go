package web

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/log"
	"github.com/goods/httpbuf"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
)

var store sessions.Store

type HandlerAdder func(string, Handle) *mux.Route
type Handle func(http.ResponseWriter, *http.Request) error

var requestLog = log.New("request")

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request) error
}

func (h Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	requestLog.Info("%v %v %v", &r, r.Method, r.RequestURI)

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
	fmt.Println(err)
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
