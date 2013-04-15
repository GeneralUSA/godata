package subjects

import (
	"database/sql"
	"fmt"
	"github.com/YouthBuild-USA/godata/questions"
	"github.com/YouthBuild-USA/godata/subjects/types"
	"github.com/YouthBuild-USA/godata/web"
	"github.com/coopernurse/gorp"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var dbMap *gorp.DbMap

func Initialize(db *gorp.DbMap) {
	dbMap = db
	types.Initialize(db)
	dbMap.AddTableWithName(Subject{}, "subject").SetKeys(true, "Id")
}

type Subject struct {
	Id            int           `db:"id"`
	SubjectTypeId int           `db:"subject_type_id"`
	Name          string        `db:"name"`
	SortName      string        `db:"sort_name"`
	RootSubjectId sql.NullInt64 `db:"root_subject_id"`
	ParentId      sql.NullInt64 `db:"parent_id"`
}

// Load loads a single subject by Id
func Load(id int) (*Subject, error) {
	item, err := dbMap.Get(Subject{}, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, fmt.Errorf("no subject with id %v", id)
	}
	return item.(*Subject), nil
}

func (s Subject) Questions() ([]*questions.Question, error) {
	return questions.ForSubjectType(s.SubjectTypeId)
}

func (s Subject) SaveQuestions(r *http.Request, subject *Subject, form url.Values) error {
	// Parse the form values into a map for each question
	questionData := parseQuestionForm(form)
	allQuestions, _ := s.Questions()

	saver := questions.DataSaver{}

	for _, question := range allQuestions {
		data := questionData[int(question.Id)]
		value, err := question.QuestionType().Parse(data)
		if err != nil {
			web.FlashError(r, fmt.Sprintf("Could not save %v: %v", question.Name, err))
			continue
		}
		ToSave := questions.Data{subject.Id, question, value}
		err = saver.SaveInstant(ToSave)
		if err != nil {
			web.FlashError(r, err.Error())
		} else {
			web.FlashInfo(r, fmt.Sprint(question.Name, value.Interface()))
		}
	}

	return nil
}

// parseQyestionForm parses a url.Values object into a url.Values object for each
// question.  Field names are prefixed by Question.<id>.
func parseQuestionForm(form url.Values) map[int]url.Values {
	questionData := make(map[int]url.Values)
	exp := regexp.MustCompile(`^Question\.(?P<id>[0-9]+)(\.(?P<name>[^\s]+))?$`)

	for name, value := range form {
		match := exp.FindStringSubmatch(name)
		if match == nil {
			continue
		}
		questionId, err := strconv.ParseInt(match[1], 10, 32)
		if err != nil {
			continue
		}
		if _, ok := questionData[int(questionId)]; !ok {
			questionData[int(questionId)] = url.Values{}
		}
		fieldName := match[3]
		questionData[int(questionId)][fieldName] = value
	}
	return questionData
}
