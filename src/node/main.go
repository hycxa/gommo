package main

import (
	"ext"
	"flag"
	"god"
	"os"
)

var (
	showHelp         bool
	noshell          bool
	nodeListenString string
)

//var clientListenString = flag.String("127.0.0.1:1119", "listen", "")
func init() {
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.BoolVar(&noshell, "noshell", false, "noshell")
	flag.StringVar(&nodeListenString, "listen", "127.0.0.1:1119", "listen string, such as 127.0.0.1:1119")
	flag.Parse()
}

func main() {
	if showHelp {
		flag.PrintDefaults()
		return
	}

	ext.TraceSwitch = true
	nodeAcceptor := god.NewWorker(god.NewAcceptor(nodeListenString, god.NewNodeAgent))
	//clientAcceptor := god.NewWorker(god.NewAcceptor(&net.TCPAddr{}, god.NewClientAgent))

	if noshell {
		c := make(chan os.Signal)
		<-c
	} else {
		god.Console().Run()
	}
	//clientAcceptor.Stop()
	nodeAcceptor.Stop()
}
