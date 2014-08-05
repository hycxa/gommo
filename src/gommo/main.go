package main

import (
	"flag"
	"god"
	"os"
)

var (
	showHelp     bool
	noshell      bool
	listenString string
)

//var clientListenString = flag.String("127.0.0.1:1119", "listen", "")
func init() {
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.BoolVar(&noshell, "noshell", false, "noshell")
	flag.StringVar(&listenString, "listen", "127.0.0.1:1119", "listen string, such as 127.0.0.1:1119")
	flag.Parse()
}

func main() {
	if showHelp {
		flag.PrintDefaults()
		return
	}

	//ext.TraceSwitch = true
	god.StartNode(listenString)
	connector := god.NewConnector(NewClientAgent)
	//clientAcceptor := god.NewWorker(god.NewAcceptor(&net.TCPAddr{}, god.NewClientAgent))

	if noshell {
		c := make(chan os.Signal)
		<-c
		god.Quit()
	} else {
		god.Console().Run()
		connector.Stop()
	}
	//clientAcceptor.Stop()
}
