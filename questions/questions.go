package questions

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/forms"
	"github.com/YouthBuild-USA/godata/questions/types"
	subjectTypes "github.com/YouthBuild-USA/godata/subjects/types"
	"github.com/YouthBuild-USA/godata/templates"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func AddRoutes(adder web.HandlerAdder) {
	types.AddRoutes(adder)
	adder("/admin/questions", questionList).Methods("GET")
	adder("/admin/questions/add/{questionType}", addQuestionForm).Methods("GET")
	adder("/admin/questions/add/{questionType}", addQuestion).Methods("POST")
	adder("/admin/questions/{id:[0-9]+}", showQuestion).Methods("GET")
	adder("/admin/questions/{id:[0-9]+}/edit", editQuestionForm).Methods("GET")
	adder("/admin/questions/{id:[0-9]+}/edit", editQuestion).Methods("POST")
}

var (
	questionListTemplate = templates.OneColumn().Add("questions/list")
	addQuestionTemplate  = templates.OneColumn().Add("questions/add", "forms/horizontal")
	viewQuestionTemplate = templates.OneColumn().Add("questions/view")
)

func showQuestion(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	questionId, _ := strconv.ParseInt(vars["id"], 10, 64)
	question, err := LoadQuestion(int(questionId))
	if err != nil {
		return err
	}
	return viewQuestionTemplate.Render(w, r, question.Name, question)
}

func questionList(w http.ResponseWriter, r *http.Request) error {
	questions, err := LoadQuestions()
	if err != nil {
		return err
	}

	return questionListTemplate.Render(w, r, "Questions", questions)
}

func editQuestionForm(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	questionId, _ := strconv.ParseInt(vars["id"], 10, 64)

	question, err := LoadQuestion(int(questionId))
	if err != nil {
		return err
	}

	var form *questionForm
	var ok bool
	if form, ok = context.Get(r, "editQuestionForm").(*questionForm); !ok {
		form = NewQuestionForm(r, question)
	}

	data := struct {
		Action string
		Form   *questionForm
	}{r.RequestURI, form}

	addQuestionTemplate.Render(w, r, "Edit Question", data)

	return nil
}

func editQuestion(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	questionId, _ := strconv.ParseInt(vars["id"], 10, 64)
	question, err := LoadQuestion(int(questionId))
	if err != nil {
		return err
	}

	r.ParseForm()

	form := NewQuestionForm(r, question)
	fmt.Println(forms.DecodeForm(form, r.Form))

	errors := forms.Validate(form)
	if len(errors) > 0 {
		context.Set(r, "editQuestionForm", form)
		editQuestionForm(w, r)
		return nil
	}

	form.scan(question)
	err = question.Save()
	if err != nil {
		return err
	}

	err = question.SetSubjectTypes(form.subjectTypeIds())
	if err != nil {
		return err
	}

	web.FlashSuccess(r, fmt.Sprintf("\"%v\" has been saved", question.Name))
	http.Redirect(w, r, fmt.Sprintf("/admin/questions/%v", question.Id), http.StatusFound)
	return nil
}

func addQuestionForm(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	questionType, ok := types.QuestionTypes[vars["questionType"]]
	if !ok {
		return fmt.Errorf("Question Type %v does not exist", vars["questionType"])
	}

	question := &Question{Type: questionType.Name()}

	var form *questionForm
	if form, ok = context.Get(r, "newQuestionForm").(*questionForm); !ok {
		form = NewQuestionForm(r, question)
	}

	data := struct {
		Action string
		Form   *questionForm
	}{r.RequestURI, form}

	return addQuestionTemplate.Render(w, r, fmt.Sprintf("Add New %v Question", question.QuestionType().Label()), data)
}

func addQuestion(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	questionType, ok := types.QuestionTypes[vars["questionType"]]
	if !ok {
		return fmt.Errorf("Question Type %v does not exist", vars["questionType"])
	}

	question := &Question{Type: questionType.Name()}

	var form = NewQuestionForm(r, question)
	r.ParseForm()
	err := forms.DecodeForm(form, r.Form)

	if err != nil {
		web.FlashWarning(r, err.Error())
	}

	if !web.ValidateCSRF(r, form.CSRF) {
		return fmt.Errorf("Invalid token")
	}

	validationErrors := forms.Validate(form)

	if len(validationErrors) == 0 {

		q := new(Question)
		q.Type = questionType.Name()
		form.scan(q)
		q.Save()

		q.SetSubjectTypes(form.subjectTypeIds())

		web.FlashSuccess(r, fmt.Sprintf("\"%v\" created successfully.", q.Name))

		http.Redirect(w, r, "/admin/questions", http.StatusFound)
		return nil
	} else {
		context.Set(r, "newQuestionForm", form)
		return addQuestionForm(w, r)
	}
	return nil
}

type questionForm struct {
	QuestionId   int64
	CSRF         string         `schema:"csrf_token"`
	Name         *forms.TextBox `schema:"QuestionName"`
	Description  *forms.Textarea
	Help         *forms.Textarea
	Text         *forms.Textarea
	Enabled      *forms.SingleCheckbox
	SubjectTypes *forms.Checkboxes
}

func NewQuestionForm(r *http.Request, question *Question) *questionForm {
	form := new(questionForm)

	form.QuestionId = question.Id

	form.Name = forms.NewTextbox("Question Name", question.Name)
	form.Name.Help = "The unique name of the question, for admin use."
	form.Name.AddValidation("Question Name is required", forms.ValidatePresence)
	form.Name.AddValidation("Question Name must be unique.", uniqueQuestionName)

	form.Description = forms.NewTextarea("Description", question.Description)
	form.Description.Help = "The description of this question, for admin use."
	form.Description.AddValidation("Description is required", forms.ValidatePresence)

	form.Text = forms.NewTextarea("Question Text", question.Text)
	form.Text.Help = "The text of the question displayed to the user"
	form.Text.AddValidation("Question Text is required", forms.ValidatePresence)

	form.Help = forms.NewTextarea("Help Text", question.Help)
	form.Help.Help = "Help text for the question, to assist the user when answering"

	form.Enabled = forms.NewCheckbox("Enabled", question.Enabled)

	types, _ := question.SubjectTypes()
	selectedTypes := make(map[int]bool, len(types))
	for _, t := range types {
		selectedTypes[t.Id] = true
	}

	form.SubjectTypes = forms.NewCheckboxes("Attach To")
	for _, t := range subjectTypes.SubjectTypeCache(r) {
		form.SubjectTypes.AddOption(strconv.Itoa(t.Id), t.Label, selectedTypes[t.Id])
	}

	fmt.Println(form.SubjectTypes.Options)

	forms.PrepareForm(form)

	return form
}

func (form *questionForm) scan(question *Question) {
	question.Name = form.Name.Value
	question.Description = form.Description.Value
	question.Text = form.Text.Value
	question.Help = form.Help.Value
	question.Enabled = form.Enabled.Checked
}

func (form *questionForm) subjectTypeIds() []int64 {
	ids := make([]int64, 0)
	for _, option := range form.SubjectTypes.Options {
		if option.Checked {
			id, _ := strconv.ParseInt(option.Key, 10, 32)
			ids = append(ids, id)
		}
	}
	fmt.Println(ids)
	return ids
}

// uniqueQuestionName is a form validator that checks a function's name is unique.
func uniqueQuestionName(form interface{}, nameField interface{}) bool {
	f := form.(*questionForm)
	field := nameField.(*forms.TextBox)
	row := DB.QueryRow("SELECT COUNT(1) FROM question WHERE Name = $1 and Id <> $2", field.Value, f.QuestionId)
	var count int
	err := row.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil && count == 0
}
