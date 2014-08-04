package god

var nodes = make(map[PID]NodeSender)

func init() {
	Console().RegCmdFunc("nodes", listNodes)
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

func listNodes([]string) interface{} {
	return nodes
}
