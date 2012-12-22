// Package template provides layouts for the data system

package templates

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/web"
	"html/template"
	"io"
	"net/http"
)

type Page struct {
	Title   string
	Body    interface{}
	Flashes []interface{}
}

var AssetPath = "assets"

var oneColumnBase = template.Must(template.ParseFiles(AssetPath + "/templates/base.html"))

type Layout struct {
	t *template.Template
}

func OneColumn() *Layout {
	return &Layout{
		t: template.Must(oneColumnBase.Clone()),
	}
}

func (l Layout) Render(w io.Writer, r *http.Request, page Page) error {
	page.Flashes = web.Flashes(r)
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
