package god

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"ext"
	"hash"
	"io"
	"net"
	"proto"
	//"time"
)

const (
	TCP_TIMEOUT = 12000
)

type NodeID hash.Hash

type NodeInfo struct {
	Name    string
	Network string
	String  string
}

type RemoteNode struct {
	NodeInfo
	net.Conn
}

type RemoteNodeMap map[string]*RemoteNode

type Node struct {
	NodeInfo
	net.Listener
	newRemote chan *RemoteNode
	delConn   chan *string
	connected RemoteNodeMap
}

var (
	nodeTrace = ext.Trace(true)
)

func NewNode(name, network, address string) *Node {
	var err error
	self := new(Node)
	self.Listener, err = net.Listen(network, address)
	if err != nil {
		ext.Errorf(err.Error())
		return nil
	}

	self.Name = name
	self.NodeInfo.Network = self.Listener.Addr().Network()
	self.NodeInfo.String = self.Listener.Addr().String()
	self.newRemote = make(chan *RemoteNode, 16)
	self.delConn = make(chan *string, 16)
	self.connected = make(map[string]*RemoteNode)
	go self.accept()
	go self.updateConnect()
	return self
}

func syncNodeInfo(conn net.Conn, nodeInfo NodeInfo) *RemoteNode {
	defer nodeTrace.UT(nodeTrace.T("Node::syncNodeInfo\t%s\tto\t%s", nodeInfo.Name, conn.RemoteAddr().String()))

	var b, wb bytes.Buffer
	var err error

	enc := gob.NewEncoder(&b)
	ext.AssertE(enc.Encode(nodeInfo))

	ext.AssertE(binary.Write(&wb, BYTE_ORDER, uint16(len(b.Bytes()))))
	_, err = conn.Write(wb.Bytes())
	ext.AssertE(err)

	_, err = conn.Write(b.Bytes())
	ext.AssertE(err)

	header := make([]byte, 2)
	_, err = io.ReadFull(conn, header)
	ext.AssertE(err)

	data := make([]byte, BYTE_ORDER.Uint16(header))
	_, err = io.ReadFull(conn, data)
	ext.AssertE(err)

	var r RemoteNode
	r.Conn = conn

	rb := bytes.NewBuffer(data)
	dec := gob.NewDecoder(rb)
	ext.AssertE(dec.Decode(&(r.NodeInfo)))
	return &r
}

func (self *Node) Dial(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return ext.LogError(err)
	}

	r := syncNodeInfo(conn, self.NodeInfo)
	self.connected[r.Name] = r

	return nil
}

func (self *Node) DialNode(target *Node) (err error) {
	return self.Dial(target.NodeInfo.Network, target.NodeInfo.String)
}

func (self *Node) accept() {
	for {
		conn, err := self.Accept()
		if err != nil {
			ext.LogError(err)
			return
		}
		self.newRemote <- syncNodeInfo(conn, self.NodeInfo)
		go self.dealOneCon(conn)
	}
}

func (self *Node) dealOneCon(conn net.Conn) {
	header := make([]byte, 2)
	//var connName string

	defer func() {
		//conn.Close()
		//self.delConn <- &connName
	}()

	for {
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(conn, header)
		ext.AssertE(err)

		data := make([]byte, BYTE_ORDER.Uint16(header))
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err = io.ReadFull(conn, data)
		ext.AssertE(err)

		b := bytes.NewBuffer(data)
		ok, msg := proto.DecodeMsg(b)
		if ok {
			_ = msg //消息处理
		}
	}
}

func (self *Node) updateConnect() {
	for {
		select {
		case newRemote, ok := <-self.newRemote:
			if ok {
				self.connected[newRemote.Name] = newRemote
			}
		case delName, ok := <-self.delConn:
			if !ok {

			}
			delete(self.connected, *delName)
		}
	}
}

func (self *Node) Connected() RemoteNodeMap {
	return self.connected
}
