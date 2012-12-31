package subjects

import (
	"database/sql"
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var rootSubjectTemplate = templates.TwoColumn().Add("subjects/rootSubject")
var subjectTemplate = templates.TwoColumn().Add("subjects/subject")

type Subject struct {
	Id            int
	SubjectTypeId int
	Name          string
	SortName      string
	Slug          sql.NullString
	Root          sql.NullInt64
}

func subjectRoutes(handleAdder web.HandlerAdder) error {
	handleAdder("/s/{root:[a-z0-9-]*}", rootSubject).Methods("GET").Name("rootSubject")
	handleAdder("/s/{root:[a-z0-9-]*}/{subjectType}/{subjectId:[0-9]+}", subjectPage).Methods("GET").Name("subject")
	return nil
}

func rootSubject(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	rootSubject, err := LoadRootSlug(vars["root"])
	if err != nil {
		return web.WebError{err, err.Error(), http.StatusNotFound}
	}

	if rootSubject.Root.Valid {
		return web.WebError{nil, "Not found", http.StatusNotFound}
	}

	return rootSubjectTemplate.Render(w, r, rootSubject.Name, rootSubject)

}

func (s Subject) String() string {
	return s.Name
}

func subjectPage(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)

	rootSubject, err := LoadRootSlug(vars["root"])
	if err != nil {
		return web.WebError{err, err.Error(), http.StatusNotFound}
	}

	subjectId, _ := strconv.Atoi(vars["subjectId"])
	subject, err := Load(subjectId)
	if err != nil {
		return web.WebError{err, err.Error(), http.StatusNotFound}
	}

	if subject.Root.Int64 != int64(rootSubject.Id) {
		return web.WebError{nil, "subject not for site", http.StatusNotFound}
	}

	return subjectTemplate.Render(w, r, subject.Name, subject)
}

func LoadRootSlug(slug string) (*Subject, error) {
	row := db.QueryRow("SELECT id, subject_type_id, name, sort_name, root_subject_id FROM subject WHERE slug = $1", slug)
	var subject Subject
	err := row.Scan(&subject.Id, &subject.SubjectTypeId, &subject.Name, &subject.SortName, &subject.Root)
	return &subject, err
}

func Load(id int) (*Subject, error) {
	row := db.QueryRow("SELECT id, subject_type_id, name, sort_name, root_subject_id FROM subject WHERE id = $1", id)
	var subject Subject
	err := row.Scan(&subject.Id, &subject.SubjectTypeId, &subject.Name, &subject.SortName, &subject.Root)
	return &subject, err
}

func (subject Subject) GetChildSubjects() ([]*Subject, error) {
	subjects := make([]*Subject, 0)
	rows, err := db.Query("SELECT id, subject_type_id, name, sort_name, root_subject_id FROM subject WHERE root_subject_id = $1", subject.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		s := new(Subject)
		err = rows.Scan(&s.Id, &s.SubjectTypeId, &s.Name, &s.SortName, &s.Root)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}
