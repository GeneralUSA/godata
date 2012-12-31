package subjects

import (
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"net/http"
)

var (
	subjecTypeIndexTemplate = templates.OneColumn().Add("subjects/typeList")
)

func subjectTypeRoutes(handleAdder web.HandlerAdder) error {
	handleAdder("/admin/subject-types", subjectTypeIndex).Methods("GET")
	return nil
}

func subjectTypeIndex(w http.ResponseWriter, r *http.Request) error {

	types, err := LoadSubjectTypes()
	if err != nil {
		return err
	}

	p := templates.Page{}
	p.Title = "Subect Types"
	p.Body = types

	subjecTypeIndexTemplate.Render(w, r, p)
	return nil
}

type SubjectType struct {
	Id   int
	Name string
}

func LoadSubjectTypes() ([]*SubjectType, error) {
	rows, err := db.Query("SELECT id, name FROM subject_type")
	if err != nil {
		return nil, err
	}

	// TODO: Figure out how to initialize the capacity to the right number...
	types := make([]*SubjectType, 0, 5)

	for rows.Next() {
		t := SubjectType{}
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			return nil, err
		}
		types = append(types, &t)
	}
	return types, nil
}
