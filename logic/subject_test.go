package logic

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSubject(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	s, _ := LoadSubject(1)
	fmt.Println(s.Name())
	s.SetName(fmt.Sprintf("Name %v", rand.Intn(1000)))
	fmt.Println(s.Name())
	s.Save()

	s2, _ := CreateSubject("student")
	fmt.Println(s2.Id())
	s2.SetName("New Student")
	s2.SetSortName("Student, New")
	s2.Save()
	fmt.Println(s2.Id())
}
