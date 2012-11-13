package datastore

import (
	"fmt"
	"testing"
)

func printSubjectFromDB(id int) {
	s, err := LoadSubject(id)
	if err != nil {
		fmt.Println("Error loading: ", err)
	}
	fmt.Println(s)
}

func TestDatastore(t *testing.T) {
	s := CreateSubject("student", "New Subject", "Subject, New")
	s.Save()
	printSubjectFromDB(s.Id)

	s.Data.Set("test key", "a data value")
	s.Save()
	printSubjectFromDB(s.Id)

	s.Data.Set("test key", "a changed value")
	s.Save()
	printSubjectFromDB(s.Id)

	s.Data.Set("new key", "a new value")
	s.Save()
	printSubjectFromDB(s.Id)

	printSubjectFromDB(s.Id)
}
