package web

import (
	"github.com/gorilla/context"
	"net/http"
)

func Cache(r *http.Request, getter func() (interface{}, error)) (interface{}, error) {
	cached := context.Get(r, getter)
	if cached != nil {
		return cached, nil
	}

	value, err := getter()
	if err != nil {
		return value, err
	}
	context.Set(r, getter, value)
	return value, nil
}
