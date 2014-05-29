package god

import (
	"testing"
	"time"
)

func sleepALittle() {
	time.Sleep(1000 * time.Millisecond)
}

func TestNewNode1(t *testing.T) {
	TestServer(t)
	NodeInit("n1", "tcp", "127.0.0.1", 8001)

	for {
		sleepALittle()
	}
}

