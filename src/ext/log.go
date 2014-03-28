package ext

import (
	// "bytes"
	"fmt"
	"log"
	"os"
	// "runtime"
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
		tl.Printf("BGN\t[%s]\n", s)
		return s
	}
	return f
}

func (t Trace) UT(s string) {
	if t {
		tl.Printf("END\t[%s]\n", s)
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

func DebugStackf(format string, v ...interface{}) {
	// Debugf(format, v...)
	// b := bytes.NewBuffer()
	// runtime.Stack(b.Bytes(), false)
	// dl.Println(b.Bytes())
}

func DebugError(err error) error {
	DebugStackf(err.Error())
	return err
}
func Errorf(format string, v ...interface{}) {
	if el == nil {
		el = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
	}
	el.Panicf(format, v...)
}

func Error(err error) error {
	Errorf(err.Error())
	return err
}
