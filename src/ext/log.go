package ext

import (
	"fmt"
	"log"
	"os"
)

type Trace bool

const (
	DEBUG = 2
	ERROR = 3
	FATAL = 4
)

var (
	tl          *log.Logger = log.New(os.Stdout, "[TRACE] ", log.LstdFlags)
	dl          *log.Logger
	el          *log.Logger
	TraceSwitch Trace
)

func init() {
}

func (t Trace) T(f string, v ...interface{}) string {
	if t {
		s := fmt.Sprintf(f, v...)
		tl.Printf("->\t[%s]\n", s)
		return s
	}
	return f
}

func (t Trace) UT(s string) {
	if t {
		tl.Printf("<-\t[%s]\n", s)
	}
}

func T(f string, v ...interface{}) string {
	return TraceSwitch.T(f, v...)
}

func UT(s string) {
	TraceSwitch.UT(s)
}

func Debugf(format string, v ...interface{}) {
	if dl == nil {
		dl = log.New(os.Stdout, "[DEBUG] ", log.LstdFlags)
	}
	dl.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	if el == nil {
		el = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	}
	el.Printf(format, v...)
	el.Print(Stack())
}

func LogError(err error) error {
	Errorf(err.Error())
	return err
}
