package main

import (
	"god"
)

func main() {
	acceptor := god.NewWorker(&god.Acceptor{})
	// node_manager, id := god.NewWorker(&god.NodeManager{})
	// worker_manager, id := god.NewWorker(&god.WorkerManager{})
	console := &god.Console{}

	console.Run()
	acceptor.Stop()
}
