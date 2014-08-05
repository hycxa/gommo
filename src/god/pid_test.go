package god

import (
	"ext"
	"testing"
)

func TestGenerate(t *testing.T) {
	count := 10
	pids := make(map[PID]bool)
	for i := 0; i < count; i++ {
		pids[GeneratePID()] = true
	}
	ext.AssertT(t, count == len(pids), "pid repeated")
}
