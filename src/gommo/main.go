package main

import (
	"ext"
	"flag"
	"god"
	"os"
)

var (
	trace        bool
	showHelp     bool
	noshell      bool
	listenString string
	agentString  string
)

func init() {
	flag.BoolVar(&trace, "trace", false, "show help")
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.BoolVar(&noshell, "noshell", false, "noshell")
	flag.StringVar(&listenString, "listen", "127.0.0.1:1119", "listen string, such as 127.0.0.1:1119")
	flag.StringVar(&agentString, "agent", "127.0.0.1:3724", "agent listen address")
	flag.Parse()
}

func main() {
	if showHelp {
		flag.PrintDefaults()
		return
	}

	if trace {
		ext.TraceSwitch = true
	}
	god.StartNode(listenString)
	acceptor := god.NewAcceptor(agentString, NewClientAgent)

	if noshell {
		c := make(chan os.Signal)
		<-c
		god.Quit()
	} else {
		god.Console().Run()
		acceptor.Stop()
	}
}
