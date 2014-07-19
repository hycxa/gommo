package god

import (
	"ext"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"sync"
	"testing"
)

func TestStartDaemon(t *testing.T) {
	runtime.GOMAXPROCS(4)
	ext.AssertT(t, runtime.GOMAXPROCS(4) == 4, "GOMAXPROCS must > 1")
	var nd Daemon
	var once sync.Once

	for i := 0; i < 100; i++ {
		go func() {
			d := StartDaemon()
			once.Do(
				func() {
					nd = d
				})
			fmt.Printf("%+v\t%+v\t%+v\t%+v\t%+v\n", i, &nd, nd, &d, d)
			ext.AssertT(t, reflect.DeepEqual(nd, d), "node daemon must be singleton.")
		}()
	}
	ext.AssertT(t, reflect.DeepEqual(nd, StartDaemon()), "node daemon must be singleton.")
}

func TestNodeDaemonListener(t *testing.T) {
	nd := StartDaemon()
	fmt.Println(nd.Network(), nd.String())
	conn, err := net.Dial(nd.Network(), nd.String())
	ext.TestingAssert(t, conn != nil && err == nil, err)
}
