package postgres

import (
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type Subject struct {
	// Data base fields
	id          int
	subjectType string
	name        string
	sortName    string
	data        *HStore

	// Other stuff
	isNew bool
}

func (s Subject) Id() int {
	return s.id
}

func (s Subject) Name() string {
	return s.name
}
func (s *Subject) SetName(name string) {
	s.name = name
}
func (s Subject) SortName() string {
	return s.sortName
}
func (s *Subject) SetSortName(sortName string) {
	s.sortName = sortName
}
func (s Subject) SubjectType() string {
	return s.subjectType
}

func newSubject() *Subject {
	return &Subject{data: NewHStore()}
}

func (pg PostgresProvider) CreateSubject(subjectType, name, sortName string) datastore.Subject {
	s := Subject{
		subjectType: subjectType,
		name:        name,
		sortName:    sortName,
		isNew:       true,
	}
	s.data = NewHStore()
	return (datastore.Subject)(&s)
}

func (pg PostgresProvider) LoadSubject(id int) (datastore.Subject, error) {
	s := newSubject()

	row := pg.statements.subject.load.QueryRow(id)
	err := row.Scan(&s.id, &s.subjectType, &s.name, &s.sortName, &s.data)
	if err != nil {
		return nil, err
	}

	return datastore.Subject(s), nil
}

func (s *Subject) Save() error {
	if s.isNew {
		return s.insert()
	} else {
		return s.update()
	}
	return nil
}

func (s *Subject) Delete() error {
	_, err := provider.statements.subject.delete.Exec(s.id)
	return err
}

func (s *Subject) insert() error {
	row := provider.statements.subject.create.QueryRow(s.subjectType, s.name, s.sortName, s.data)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	s.id = id
	s.isNew = false
	return nil
}

func (s Subject) update() error {
	_, err := provider.statements.subject.update.Exec(s.id, s.name, s.sortName)
	s.data.update(s.id,
		provider.statements.subject.updateData,
		provider.statements.subject.deleteData,
	)
	return err
}
