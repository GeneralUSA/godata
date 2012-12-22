package godata

import (
	"database/sql"
	"github.com/YouthBuild-USA/godata/log"
	_ "github.com/bmizerany/pq"
	"os"
	"testing"
)

func TestModule(t *testing.T) {
	log.Add(log.DEBUG, os.Stdout)
	db, _ := sql.Open("postgres", "user=andrew password=root dbname=andrew")

	SetDatabase(db)
	Start()
}
