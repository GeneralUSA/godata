package types

import (
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"net/http"
	"net/url"
	"reflect"
)

var (
	typeListTemplate = templates.OneColumn().Add("questions/types/list")
)

type QuestionType interface {
	Name() string
	Label() string
	Description() string
	Template() string
	FieldType() FieldType
	Parse(form url.Values) (reflect.Value, error)
}

var QuestionTypes = make(map[string]QuestionType)

func Register(questionType QuestionType) {
	if _, ok := QuestionTypes[questionType.Name()]; !ok {
		QuestionTypes[questionType.Name()] = questionType
	}
}

func AddRoutes(adder web.HandlerAdder) {
	adder("/admin/questions/types", subjectTypeList).Methods("GET")
}

func subjectTypeList(w http.ResponseWriter, r *http.Request) error {

	typeListTemplate.Render(w, r, "Question Types", QuestionTypes)

	return nil
}
