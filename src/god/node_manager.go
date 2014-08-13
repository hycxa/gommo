package god

import (
	"ext"
)

type nodeInfo struct {
	Cookie string
	ID     NID
}

const (
	RELEASE = false
)

var (
	nodes         = make(map[NID]interface{})
	nodeConnector Connector
	nodeAcceptor  Stopper
	myInfo        nodeInfo
)

func regCmd() {
	Console().RegCmdFunc("dial", dial)
	Console().RegCmdFunc("nodes", listNodes)
	Console().RegCmdFunc("q", quit)
	Console().RegCmdFunc("quit", quit)
}

func StartNode(listenStr string) {
	nodeAcceptor = NewAcceptor(listenStr, NewNodeAgent)
	nodeConnector = NewConnector(NewNodeAgent)
	myInfo = nodeInfo{Cookie: "THIS_IS_A_COOKIE", ID: GenerateNID()}
	ext.LogInfo("MY_INFO\t%v", myInfo)
	regCmd()
}

func MyInfo() nodeInfo {
	return myInfo
}

func Quit() {
	if nodeAcceptor != nil {
		nodeAcceptor.Stop()
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

func AddNode(id NID, node interface{}) {
	nodes[id] = node
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
