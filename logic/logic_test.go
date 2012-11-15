package logic

import (
	"fmt"
	"github.com/YouthBuild-USA/ybdata/datastore"
	"github.com/YouthBuild-USA/ybdata/datastore/postgres"
	"math/rand"
	"testing"
	"time"
)

func TestSubject(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	p := postgres.NewPostgresProvidier("localhost", "andrew", "root", "andrew")
	datastore.Register(p)

	s, _ := LoadSubject(1)
	fmt.Println(s.Name())
	s.SetName(fmt.Sprintf("Name %v", rand.Intn(1000)))
	s.Save()
	fmt.Println(s.Name())

}
