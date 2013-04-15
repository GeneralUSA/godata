package types

import (
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/context"
	"net/http"
)

var dbMap *gorp.DbMap

func Initialize(db *gorp.DbMap) {
	dbMap = db
	dbMap.AddTableWithName(SubjectType{}, "subject_type").SetKeys(true, "Id")
}

var (
	subjecTypeIndexTemplate = templates.OneColumn().Add("subjects/typeList")
)

func AddRoutes(handleAdder web.HandlerAdder) error {
	handleAdder("/admin/subject-types", subjectTypeIndex).Methods("GET")
	return nil
}

func subjectTypeIndex(w http.ResponseWriter, r *http.Request) error {

	types, err := LoadSubjectTypes()
	if err != nil {
		return err
	}

	subjecTypeIndexTemplate.Render(w, r, "Subject Types", types)
	return nil
}

type SubjectType struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Label string `db:"label"`
}

func LoadSubjectTypes() ([]*SubjectType, error) {
	items, err := dbMap.Select(SubjectType{}, "SELECT id, name, label FROM subject_type")
	if err != nil {
		return nil, err
	}

	types := make([]*SubjectType, 0)
	for _, i := range items {
		types = append(types, i.(*SubjectType))
	}
	return types, nil
}

type contextKeys int

const (
	SubjectTypeKey = contextKeys(iota)
)

func SubjectTypeCache(r *http.Request) []*SubjectType {
	val := context.Get(r, SubjectTypeKey)
	if val != nil {
		return val.([]*SubjectType)
	}

	types, _ := LoadSubjectTypes()
	context.Set(r, SubjectTypeKey, types)
	return types
}
