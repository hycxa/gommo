package ext

import "fmt"

type Debug bool

func (d Debug) Printf(s string, a ...interface{}) {
	if d {
		fmt.Printf(s, a...)
	}
}

func DebugLog() {

}
