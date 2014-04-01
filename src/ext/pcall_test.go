package ext

import (
	"testing"
)

func TestPCall(t *testing.T) {
	PCall(func() {
		a := 0
		b := 1000 / a
		_ = b
	})
}
