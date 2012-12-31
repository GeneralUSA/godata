package godata

import (
	"github.com/YouthBuild-USA/godata/log"
	"os"
	"testing"
)

func TestModule(t *testing.T) {
	log.Add(log.DEBUG, os.Stdout)
	Start("config.cfg")
}
