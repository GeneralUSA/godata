package postgres

import (
	"fmt"
)

type Subject struct {
	// Data base fields
	Id       int
	Type     string
	Name     string
	SortName string
	Data     *HStore

	// Other stuff
	IsNew bool
}

var subjectStatements *SubjectStatements

func init() {
	db := getConnection()
	subjectStatements = newSubjectStatements(db)
}

func newSubject() *Subject {
	return &Subject{Data: NewHStore()}
}

func CreateSubject(typ string, name string, sortName string) *Subject {
	s := Subject{
		Type:     typ,
		Name:     name,
		SortName: sortName,
		IsNew:    true,
	}
	s.Data = NewHStore()
	return &s
}

func (subject Subject) String() string {
	s := subject.Name
	for _, k := range subject.Data.Keys() {
		v, _ := subject.Data.Get(k)
		s += fmt.Sprintf("\n  %v: %v", k, v)
	}
	return s
}
func LoadSubject(id int) (s *Subject, err error) {
	s = newSubject()

	row := subjectStatements.load.QueryRow(id)
	err = row.Scan(&s.Id, &s.Type, &s.Name, &s.SortName, &s.Data)
	if err != nil {
		return
	}

	return
}

func (s *Subject) Save() error {
	if s.IsNew {
		return s.insert()
	} else {
		return s.update()
	}
	return nil
}

func (s Subject) Delete() error {
	_, err := subjectStatements.delete.Exec(s.Id)
	return err
}

func (s *Subject) insert() error {
	row := subjectStatements.create.QueryRow(s.Type, s.Name, s.SortName, s.Data)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	s.Id = id
	s.IsNew = false
	return nil
}

func (s Subject) update() error {
	_, err := subjectStatements.update.Exec(s.Id, s.Name, s.SortName)
	s.Data.update(s.Id, subjectStatements.updateData, subjectStatements.deleteData)
	return err
}
