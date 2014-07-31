package main

import (
	"god"
	"runtime"
)

func main() {
	acceptor := god.NewWorker(&god.Acceptor)
	node_manager := god.NewWorker(&god.NodeManager)
	worker_manager := god.NewWorker(&god.WorkerManager)
	console := god.NewWorker(&god.Console)
	for console.Working() {
		runtime.Gosched()
	}
}
