package ext

import (
	// "bytes"
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

func (t Trace) trace(f string) string {
	if t {
		tl.Printf(">>>> [%s]\n", f)
	}
	return f
}

func (t Trace) untrace(f string) {
	if t {
		tl.Printf("<<<< [%s]\n", f)
	}
}

func T(f string) string {
	return TraceSwitch.trace(f)
}

func UT(f string) {
	TraceSwitch.untrace(f)
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
	//el.Panicf(format, v...)
}

func Error(err error) error {
	Errorf(err.Error())
	return err
}
