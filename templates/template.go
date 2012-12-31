// Package template provides layouts for the data system

package templates

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/config"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

func init() {
	config.Register("branding", "siteName", "Go Data", "The name of the system")
}

type meta struct {
	SiteName string
}

type page struct {
	Meta    meta
	Title   string
	Data    interface{}
	User    interface{}
	Flashes []interface{}
}

func newPage(r *http.Request, title string, data interface{}) *page {
	p := &page{
		Title:   title,
		Data:    data,
		Flashes: web.Flashes(r),
		User:    context.Get(r, "user"),
		Meta: meta{
			SiteName: config.MustGet("branding", "siteName"),
		},
	}

	return p
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

func (l Layout) Render(w io.Writer, r *http.Request, title string, data interface{}) error {
	page := newPage(r, title, data)
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
