package main

import (
	"flag"
	"god"
)

var (
	showHelp      bool
	serverAddress string
	amount        int
)

//var clientListenString = flag.String("127.0.0.1:1119", "listen", "")
func init() {
	flag.BoolVar(&showHelp, "help", false, "show help")
	flag.StringVar(&serverAddress, "server", "127.0.0.1:3724", "server address")
	flag.IntVar(&amount, "amount", 100, "robot amount")
	flag.Parse()
}

func main() {
	if showHelp {
		flag.PrintDefaults()
		return
	}

	//ext.TraceSwitch = true
	connector := god.NewConnector(NewRobotAgent)
	for i := 0; i < amount; i++ {
		connector.Dial(serverAddress)
	}

	god.Console().Run()
	connector.Stop()
}
