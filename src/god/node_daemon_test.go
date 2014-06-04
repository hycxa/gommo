package god

import (
	"ext"
	"fmt"
	"net"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestStartNodeDaemon(t *testing.T) {
	runtime.GOMAXPROCS(4)
	ext.AssertT(t, runtime.GOMAXPROCS(4) == 4, "GOMAXPROCS must > 1")
	var nd Daemon

	for i := 0; i < 100; i++ {
		go func() {
			d := StartNodeDaemon()
			nd = d
			//fmt.Printf("%+v\t%+v\t%+v\t%+v\t%+v\n", i, &nd, nd, &d, d)
			ext.AssertT(t, reflect.DeepEqual(nd, d), "node daemon must be singleton.")
		}()
	}
	time.Sleep(time.Second)
}

func TestNodeDaemonListener(t *testing.T) {
	nd := StartNodeDaemon()
	fmt.Println(nd.Network(), nd.String())
	conn, err := net.Dial(nd.Network(), nd.String())
	ext.AssertT(t, conn != nil && err == nil, "listener wrong")
}
