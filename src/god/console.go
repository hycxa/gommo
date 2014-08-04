package god

/*
#include <readline/readline.h>
#include <stdlib.h>
#cgo LDFLAGS: -lreadline
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Console struct {
}

func (c *Console) Run() {
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
		fmt.Printf("%v\n", line)
	}
}
