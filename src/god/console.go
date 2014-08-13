package god

/*
#include <readline/readline.h>
#include <readline/history.h>
#include <stdlib.h>
#cgo LDFLAGS: -lreadline
*/
import "C"

import (
	"ext"
	"fmt"
	"regexp"
	"unsafe"
)

type cmdFunc func([]string) interface{}
type cmdFuncMap map[string]cmdFunc

type console struct {
	*stopper
	funcs cmdFuncMap
}

var (
	cons = &console{funcs: make(cmdFuncMap), stopper: NewStopper()}
)

func Console() *console {
	return cons
}

func (c *console) Run() {
	defer c.Stopped()

	re := regexp.MustCompile(`[\w\:\.]+`)
	prompt := C.CString("god> ")
	defer C.free(unsafe.Pointer(prompt))

	var line string

	for !c.StopRequested() {
		cline := C.readline(prompt)
		defer C.free(unsafe.Pointer(cline))
		if cline == nil {
			fmt.Printf("\n")
			break
		}

		C.add_history(cline)
		line = C.GoString(cline)
		args := re.FindAllString(line, -1)
		if len(args) > 0 {
			f := c.funcs[args[0]]
			if f != nil {
				var ret interface{}
				ext.PCall(
					func() {
						ret = f(args[1:])
					})
				ext.LogInfo("RUN_COMMAND\t%s\t%s\t%v\n", args[0], args[1:], ret)
			}
		}
	}
}

func (c *console) RegCmdFunc(cmd string, f cmdFunc) {
	c.funcs[cmd] = f
}
