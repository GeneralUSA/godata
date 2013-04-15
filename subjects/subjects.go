package subjects

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/config"
	"github.com/YouthBuild-USA/godata/questions"
	"github.com/YouthBuild-USA/godata/subjects/types"
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func init() {
	config.Register("Branding", "rootSubjectURIPrefix", "s", "The Prefix in the URI for root subjects.")
}

var rootSubjectTemplate = templates.TwoColumn().Add("subjects/rootSubject")
var subjectTemplate = templates.TwoColumn().Add("subjects/subject")

func AddRoutes(handleAdder web.HandlerAdder) error {
	types.AddRoutes(handleAdder)
	rootURI := config.MustGet("Branding", "rootSubjectURIPrefix")
	handleAdder(fmt.Sprintf("/%s/{id:[0-9]+}", rootURI), simpleSubjectPage).Methods("GET")
	handleAdder("/s/{id:[0-9]+}", saveSubject).Methods("POST")
	return nil
}

func simpleSubjectPage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subjectId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		return err
	}
	return subjectPage(w, r, int(subjectId))
}

func subjectPage(w http.ResponseWriter, r *http.Request, subjectId int) error {
	subject, err := Load(subjectId)
	if err != nil {
		return web.WebError{err, err.Error(), http.StatusNotFound}
	}

	var data struct {
		Subject   *Subject
		Questions []*questions.Question
	}

	data.Subject = subject
	data.Questions, err = subject.Questions()
	if err != nil {
		return err
	}

	return subjectTemplate.Render(w, r, subject.Name, data)
}

func saveSubject(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	subjectId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		return err
	}

	subject, err := Load(int(subjectId))
	if err != nil {
		return err
	}

	err = r.ParseForm()
	if err != nil {
		return err
	}

	subject.SaveQuestions(r, subject, r.Form)

	http.Redirect(w, r, "/s/"+strconv.Itoa(int(subjectId)), http.StatusFound)
	return nil
}
