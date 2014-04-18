package god

import (
	"testing"
	"time"
)

func sleepALittle() {
	time.Sleep(1000 * time.Millisecond)
}

func TestNewNode(t *testing.T) {
	TestServer(t)

	for {
		sleepALittle()
	}
}
