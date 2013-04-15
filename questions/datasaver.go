package questions

import (
	"fmt"
	"github.com/YouthBuild-USA/godata/global"
	"reflect"
)

type DataSaver struct {
}

type Data struct {
	SubjectId int
	Question  *Question
	Value     reflect.Value
}

func (ds DataSaver) SaveInstant(data Data) error {
	result, err := global.DbMap.Exec(fmt.Sprintf("UPDATE %v SET value = $1 WHERE question_id = $2 AND subject_id = $3", data.Question.QuestionType().FieldType().Table()),
		data.Value.Interface(),
		data.Question.Id,
		data.SubjectId,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		_, err = global.DbMap.Exec(fmt.Sprintf("INSERT INTO %v (question_id, subject_id, value, user_id, last_update) VALUES ($1, $2, $3, $4, NOW())", data.Question.QuestionType().FieldType().Table()),
			data.Question.Id,
			data.SubjectId,
			data.Value.Interface(),
			0,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
