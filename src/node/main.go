package main

import (
	"god"
	"net"
)

func main() {
	nodeAcceptor := god.NewWorker(god.NewAcceptor(&net.TCPAddr{}, &god.NodeAgentCreator{}))
	clientAcceptor := god.NewWorker(god.NewAcceptor(&net.TCPAddr{}))
	// node_manager, id := god.NewWorker(&god.NodeManager{})
	// worker_manager, id := god.NewWorker(&god.WorkerManager{})
	console := &god.Console{}

	console.Run()
	clientAcceptor.Stop()
	nodeAcceptor.Stop()
}
