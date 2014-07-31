package god

var (
	nodeManager NodeManager
)

type NodeManager struct {
	nodeSenders map[PID]NodeSender
}

func (n *NodeManager) Add(pid PID, nodeSender NodeSender) {
	n.nodeSenders[pid] = nodeSender
}
