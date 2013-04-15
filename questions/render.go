package questions

import (
	"bytes"
	"fmt"
	"github.com/YouthBuild-USA/godata/questions/types"
	"html/template"
)

var t *template.Template

var assetPath = "assets"

func init() {
	t = template.New("")

	files := make(map[string]bool)
	fileList := make([]string, 0)
	for _, t := range types.QuestionTypes {
		if !files[t.Template()] {
			files[t.Template()] = true
			fileList = append(fileList, fmt.Sprintf("%v/templates/questions/form/%v.html", assetPath, t.Template()))
		}
	}

	template.Must(t.ParseFiles(fileList...))
}

func (q *Question) Render(context interface{}) interface{} {
	buffer := new(bytes.Buffer)
	t.ExecuteTemplate(buffer, q.QuestionType().Template(), q)
	return template.HTML(buffer.String())
}
