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

func sep() {
	fmt.Print("\n----------------------\n\n")
}

func TestDatastore(t *testing.T) {
	p := NewPostgresProvidier("localhost", "andrew", "root", "andrew")
	datastore.Register(p)

	sep()

	pt, _ := p.LoadPrimitiveType(1)
	fmt.Println(pt)

	sep()

	vt, _ := p.LoadValueType(1)
	fmt.Println(vt)

	sep()

	d, _ := p.LoadDatapoint(1)
	fmt.Println(d)

	sep()

}
