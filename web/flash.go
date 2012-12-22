package web

import (
	"encoding/gob"
	"fmt"
	"net/http"
)

type flashLevel int

type FlashMessage struct {
	Level   flashLevel
	Message string
}

const (
	flashInfo flashLevel = iota
	flashSuccess
	flashWarning
	flashError
)

func (fm FlashMessage) Classes() string {
	switch fm.Level {
	case flashInfo:
		return "alert-info"
	case flashSuccess:
		return "alert-success"
	case flashError:
		return "alert-error"
	}
	return ""
}

func init() {
	gob.Register(FlashMessage{})
}

func flash(r *http.Request, level flashLevel, message string) {
	session := Session(r)
	session.AddFlash(FlashMessage{level, message})
	fmt.Println(session.Values)
}

func FlashInfo(r *http.Request, message string) {
	flash(r, flashInfo, message)
}

func FlashWarning(r *http.Request, message string) {
	flash(r, flashWarning, message)
}

func FlashSuccess(r *http.Request, message string) {
	flash(r, flashSuccess, message)
}
func FlashError(r *http.Request, message string) {
	flash(r, flashError, message)
}

func Flashes(r *http.Request) []interface{} {
	return Session(r).Flashes()
	// flashes := Session(r).Flashes()
	// messages := make([]*FlashMessage, 0, len(flashes))
	// for i := range flashes {
	//  if message, ok := flashes[i].(FlashMessage); ok {
	//    messages = append(messages, message)
	//  }
	// }
	// return messages
}
