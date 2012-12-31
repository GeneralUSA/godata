// Package template provides layouts for the data system

package templates

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type Page struct {
	Title   string
	Body    interface{}
	User    interface{}
	Flashes []interface{}
}

var Router *mux.Router

var AssetPath = "assets"

var oneColumnBase = newLayout("base")
var twoColumnBase = newLayout("base", "twoCol")

type Layout struct {
	t *template.Template
}

func newLayout(baseTemplates ...string) *template.Template {
	files := make([]string, len(baseTemplates))
	for i := range baseTemplates {
		files[i] = fmt.Sprintf("%v/templates/%v.html", AssetPath, baseTemplates[i])
	}
	funcMap := template.FuncMap{
		"url": routerUrl,
	}
	t := template.Must(template.ParseFiles(files...))
	t.Funcs(funcMap)
	return t
}

func OneColumn() *Layout {
	return &Layout{
		t: template.Must(oneColumnBase.Clone()),
	}
}

func TwoColumn() *Layout {
	return &Layout{
		t: template.Must(twoColumnBase.Clone()),
	}
}

func (l Layout) Render(w io.Writer, r *http.Request, page Page) error {
	page.Flashes = web.Flashes(r)
	page.User = context.Get(r, "user")
	return l.t.Execute(w, page)
}

func (l *Layout) Add(templates ...string) *Layout {
	files := make([]string, len(templates))
	for i := range templates {
		files[i] = fmt.Sprintf("%v/templates/%v.html", AssetPath, templates[i])
	}
	template.Must(l.t.ParseFiles(files...))
	return l
}

func routerUrl(name string, params ...string) (*url.URL, error) {
	return Router.Get(name).URLPath(params...)
}
