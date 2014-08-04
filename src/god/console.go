package god

/*
#include <readline/readline.h>
#include <stdlib.h>
#cgo LDFLAGS: -lreadline
*/
import "C"

import (
	"fmt"
	"regexp"
	"unsafe"
)

type gmCommand struct {
	command string
	args    []string
}

type Console struct {
}

func (c *Console) Run() {
	re := regexp.MustCompile(`\w+`)
	prompt := C.CString("god> ")
	defer C.free(unsafe.Pointer(prompt))
	var line string

	for line != "q" {
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
			cmd := gmCommand{args[0], args[1:]}
			Cast(0, 0, &cmd)
			fmt.Printf("%v\n", cmd)
		}
	}
}
