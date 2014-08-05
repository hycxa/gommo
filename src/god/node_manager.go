package god

import (
	"ext"
)

var (
	nodes         = make(map[PID]NodeSender)
	nodeConnector = NewConnector(NewNodeAgent)
	nodeAcceptor  Worker
)

func init() {
	Console().RegCmdFunc("dial", dial)
	Console().RegCmdFunc("nodes", listNodes)
	Console().RegCmdFunc("q", quit)
	Console().RegCmdFunc("quit", quit)
}

func Start(listenStr string) {
	nodeAcceptor = NewWorker(NewAcceptor(listenStr, NewNodeAgent))
}

func FindWorker(pid PID) Worker {
	return nil

}

func FindNodeOfWorker(pid PID) NodeSender {
	return nil
}

func Cast(source PID, target PID, message Message) {
	worker := FindWorker(target)
	if worker != nil {
		worker.Cast(source, message)
	}

	sender := FindNodeOfWorker(target)
	if sender != nil {
		sender.Cast(source, target, message)
	}
}

func AddNode(pid PID, nodeSender NodeSender) {
	nodes[pid] = nodeSender
}

func dial(args []string) interface{} {
	ext.PCall(
		func() {
			nodeConnector.Dial(args[0])
		})
	return true
}

func listNodes([]string) interface{} {
	return nodes
}

func quit([]string) interface{} {
	nodeAcceptor.Stop()
	Console().Stop()
	return true
}
