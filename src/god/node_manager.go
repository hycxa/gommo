package god

import ()

var (
	nodes         = make(map[PID]NodeSender)
	nodeConnector Connector
	nodeAcceptor  Worker
)

func init() {
	Console().RegCmdFunc("dial", dial)
	Console().RegCmdFunc("nodes", listNodes)
	Console().RegCmdFunc("q", quit)
	Console().RegCmdFunc("quit", quit)
}

func StartNode(listenStr string) {
	nodeAcceptor = NewWorker(NewAcceptor(listenStr, NewNodeAgent))
	nodeConnector = NewConnector(NewNodeAgent)
}

func Quit() {
	if nodeAcceptor != nil {
		nodeAcceptor.Stop()
	}
	if nodeConnector != nil {
		nodeConnector.Stop()
	}
	Console().Stop()
}

func Dial(address string) {
	nodeConnector.Dial(address)
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
	Dial(args[0])
	return true
}

func listNodes([]string) interface{} {
	return nodes
}

func quit([]string) interface{} {
	Quit()
	return true
}
