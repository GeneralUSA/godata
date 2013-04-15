package questions

import (
	"database/sql"
	"errors"
	"github.com/YouthBuild-USA/godata/questions/types"
	subjectTypes "github.com/YouthBuild-USA/godata/subjects/types"
	"github.com/coopernurse/gorp"
)

var DB *sql.DB
var dbMap *gorp.DbMap

func InitializeDatabase(db *sql.DB, gorpMap *gorp.DbMap) {
	DB = db
	dbMap = gorpMap
	dbMap.AddTableWithName(Question{}, "question").SetKeys(true, "Id")
	dbMap.AddTableWithName(SubjectQuestion{}, "subject_question").SetKeys(true, "Id")
}

type Question struct {
	Id           int64                       `db:"id"`
	Name         string                      `db:"name"`
	Enabled      bool                        `db:"enabled"`
	Type         string                      `db:"type"`
	Text         string                      `db:"text"`
	Description  string                      `db:"description"`
	Help         string                      `db:"help"`
	Validator    string                      `db:"validation"`
	subjectTypes []*subjectTypes.SubjectType `db:"-"`
}

func (q Question) QuestionType() types.QuestionType {
	return types.QuestionTypes[q.Type]
}

func (q *Question) Save() error {
	var err error
	if q.Id == 0 {
		err = dbMap.Insert(q)
	} else {
		_, err = dbMap.Update(q)
	}
	return err
}

func LoadQuestion(id int) (*Question, error) {
	q, err := dbMap.Get(Question{}, id)
	return q.(*Question), err
}

func LoadQuestions() ([]*Question, error) {
	s := `SELECT id, name, enabled, type, text, description, help, validation FROM question`
	q, err := dbMap.Select(Question{}, s)
	if err != nil {
		return nil, err
	}
	qs := make([]*Question, 0, len(q))
	for _, item := range q {
		qs = append(qs, item.(*Question))
	}
	return qs, nil
}

type SubjectQuestion struct {
	Id            int64 `db:"id"`
	SubjectTypeId int64 `db:"subject_type_id"`
	QuestionId    int64 `db:"question_id"`
	Weight        int64 `db:"weight"`
}

func (q Question) SetSubjectTypes(ids []int64) error {
	if q.Id == 0 {
		return errors.New("Cannot set subject types for question with no id")
	}

	// Todo(andrew) this wastes ids.
	dbMap.Exec("DELETE FROM subject_question WHERE question_id = $1", q.Id)

	sqs := make([]interface{}, 0, len(ids))

	for _, id := range ids {
		sqs = append(sqs, &SubjectQuestion{
			SubjectTypeId: id,
			QuestionId:    q.Id,
		})
	}

	return dbMap.Insert(sqs...)
}

func (q *Question) SubjectTypes() ([]*subjectTypes.SubjectType, error) {
	if q.subjectTypes != nil {
		return q.subjectTypes, nil
	}

	if q.Id == 0 {
		return nil, errors.New("Cannot load subject types for question with no id")
	}

	// TODO(andrew) load the following with gorp
	rows, err := DB.Query("SELECT subject_type.id, name, label FROM subject_type INNER JOIN subject_question ON subject_type.id = subject_type_id WHERE question_id = $1", q.Id)
	if err != nil {
		return nil, err
	}

	types := make([]*subjectTypes.SubjectType, 0)
	var t *subjectTypes.SubjectType
	for rows.Next() {
		t = new(subjectTypes.SubjectType)
		err = rows.Scan(&t.Id, &t.Name, &t.Label)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	q.subjectTypes = types
	return types, nil
}

func ForSubjectType(id int) ([]*Question, error) {
	statement := `
	SELECT q.id, q.name, q.type, q.text
	FROM question q
	INNER JOIN subject_question qs
		ON qs.question_id = q.id
	WHERE qs.subject_type_id = $1
	  AND q.enabled
	`

	items, err := dbMap.Select(Question{}, statement, id)
	if err != nil {
		return nil, err
	}

	questions := make([]*Question, 0)
	for _, item := range items {
		questions = append(questions, item.(*Question))
	}
	return questions, nil
}
