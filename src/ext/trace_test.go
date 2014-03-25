package ext

import (
	"testing"
)

func TurnTraceOff() {
	TraceSwitch = false
}

func TestTrace(t *testing.T) {
	TraceSwitch = true
	defer TurnTraceOff()
	defer UT(T("TestTrace"))
}
