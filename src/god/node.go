package god

import (
	"bytes"
	"encoding/gob"
	"ext"
	"net"
)

type NodeInfo struct {
	Name    string
	Network string
	String  string
}

type RemoteNode struct {
	NodeInfo
	net.Conn
}

type RemoteNodeMap map[string]RemoteNode

type Node struct {
	NodeInfo
	net.Listener
	c         chan NodeInfo
	connected RemoteNodeMap
}

func NewNode(name, network, address string) *Node {
	var err error
	n := new(Node)
	n.Listener, err = net.Listen(network, address)
	if err != nil {
		ext.Errorf(err.Error())
		return nil
	}

	n.Name = name
	n.NodeInfo.Network = n.Listener.Addr().Network()
	n.NodeInfo.String = n.Listener.Addr().String()
	n.connected = make(map[string]RemoteNode)
	return n
}

func syncNodeInfo(conn net.Conn, mine, remote *NodeInfo) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(mine); err != nil {
		return err
	}
	_, err := conn.Write(b.Bytes())
	return err
}

func (n *Node) Dial(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return ext.Error(err)
	}

	var r RemoteNode
	if err := syncNodeInfo(conn, &n.NodeInfo, &r.NodeInfo); err != nil {
		return ext.Error(err)
	}

	n.connected[r.Name] = r

	return nil
}

func (n *Node) DialNode(target *Node) (err error) {
	return n.Dial(target.NodeInfo.Network, target.NodeInfo.String)
}

func (n *Node) accept() {
	conn, err := n.Listener.Accept()
	if err != nil {
		ext.Error(err)
		return
	}

	var r RemoteNode
	if err := syncNodeInfo(conn, &n.NodeInfo, &r.NodeInfo); err != nil {
		ext.Error(err)
		return
	}

	n.connected[r.Name] = r
}

func (n *Node) Connected() RemoteNodeMap {
	return n.connected
}
