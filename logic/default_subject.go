package logic

import (
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type DefaultSubject struct {
	subj *datastore.Subject
}

func createDefaultSubject(subjectType string) Subject {
	s := DefaultSubject{
		datastore.CreateSubject(subjectType, "", ""),
	}
	return s
}

func loadDefaultSubject(s *datastore.Subject) Subject {
	return DefaultSubject{s}
}

func (s DefaultSubject) Id() int {
	return s.subj.Id
}

func (s DefaultSubject) Name() string {
	return s.subj.Name
}

func (s DefaultSubject) SetName(name string) {
	s.subj.Name = name
}

func (s DefaultSubject) SortName() string {
	return s.subj.SortName
}

func (s DefaultSubject) SetSortName(sortName string) {
	s.subj.SortName = sortName
}

func (s DefaultSubject) Type() string {
	return s.subj.Type
}

func (s DefaultSubject) Save() (err error) {
	err = s.subj.Save()
	return
}
