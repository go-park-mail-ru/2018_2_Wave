package logger

import (
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestLog(t *testing.T) {
	filename := "_tmp_.log"

	// test file writing
	// test file appending
	for i := 0; i < 2; i++ {
		log, err := New(Config{
			File: filename,
		})
		if err != nil {
			t.Error(err)
		}

		log.Infoln("test1")
		log.Warnln("test2")
		log.Errorln("test3")
		log.Close()
	}
	defer os.Remove(filename)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
	}
	result := strings.Split(regexp.MustCompile(`"time":"[-+:\w]+"`).
		ReplaceAllString(string(data), `"time":""`),
		"\n")

	expected := []string{`{"level":"info","msg":"test1","time":""}`,
		`{"level":"warning","msg":"test2","time":""}`,
		`{"level":"error","msg":"test3","time":""}`,
		`{"level":"info","msg":"test1","time":""}`,
		`{"level":"warning","msg":"test2","time":""}`,
		`{"level":"error","msg":"test3","time":""}`,
		``,
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected log result:\n\tExpected:\n%s\n\tTaken:\n%s\n", expected, result)
	}
}
