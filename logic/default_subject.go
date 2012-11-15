package logic

import (
	"github.com/YouthBuild-USA/ybdata/datastore"
)

type DefaultSubject struct {
	datastore.Subject
}

func createDefaultSubject(subjectType string) Subject {
	s := datastore.Provider.CreateSubject(subjectType, "", "")
	return (Subject)(s)
}

func loadDefaultSubject(s datastore.Subject) Subject {
	return s
}
