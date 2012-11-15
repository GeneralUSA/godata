package postgres

import (
	"fmt"
	"github.com/YouthBuild-USA/ybdata/datastore"
	"testing"
)

func printSubjectFromDB(id int) {
	s, err := provider.LoadSubject(id)
	if err != nil {
		fmt.Println("Error loading: ", err)
	}
	fmt.Println(s)
}

func TestDatastore(t *testing.T) {
	p := NewPostgresProvidier("localhost", "andrew", "root", "andrew")
	datastore.Register(p)

	s := p.CreateSubject("student", "New Subject", "Subject, New")
	s.Save()
	printSubjectFromDB(s.Id())
}
