package log

import (
	"bytes"
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	output := new(bytes.Buffer)
	logger := log.New(output, "", 0)
	testAspect := New("test")
	testAspect.Debug("%v", "Debug message")
	testAspect.Info("%v", "Info message")
	testAspect.Critical("%v", "Critical Message")

	expected := `Info message
Critical Message
`
	if output.String() != expected {
		t.Error("Expected output did not match")
	}
}
