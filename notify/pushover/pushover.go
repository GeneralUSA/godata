package pushover

import (
	"net/http"
	"net/url"
)

type Pushover struct {
	App   string
	Users []string
}

func New(appKey string, userKeys ...string) *Pushover {
	return &Pushover{
		App:   appKey,
		Users: userKeys,
	}
}

func (pushover Pushover) Write(p []byte) (int, error) {
	message := string(p)
	pushover.Send(message)
	return len(p), nil
}

func (pushover Pushover) Send(message string) {
	for _, key := range pushover.Users {
		send(pushover.App, key, message)
	}
}

func send(appkey, userkey, message string) {
	data := url.Values{}
	data.Add("token", appkey)
	data.Add("user", userkey)
	data.Add("message", message)
	_, err := http.PostForm("https://api.pushover.net/1/messages.json", data)
	if err != nil {
		panic(err)
	}
}
