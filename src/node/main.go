package main

import (
	"ext"
	"flag"
	"god"
	"net"
)

var species = flag.String("species", "host", "")

func main() {
	ext.TraceSwitch = true
	nodeAcceptor := god.NewWorker(god.NewAcceptor(&net.TCPAddr{}, god.NewNodeAgent))
	clientAcceptor := god.NewWorker(god.NewAcceptor(&net.TCPAddr{}, god.NewClientAgent))
	// node_manager, id := god.NewWorker(&god.NodeManager{})
	// worker_manager, id := god.NewWorker(&god.WorkerManager{})

	god.Console().Run()
	clientAcceptor.Stop()
	nodeAcceptor.Stop()
}
