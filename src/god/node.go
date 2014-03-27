package god

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"ext"
	"io"
	"net"
	"proto"
	//"time"
)

const (
	TCP_TIMEOUT = 12000
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
	newConn   chan *RemoteNode
	delConn   chan *string
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
	n.newConn = make(chan *RemoteNode)
	n.delConn = make(chan *string)
	n.connected = make(map[string]RemoteNode)
	go n.accept()
	go n.updateConnect()
	return n
}

func SendConnMsg(conn net.Conn, b *bytes.Buffer) error {
	v := b.Len()
	buf := make([]byte, 2, 2+v)
	buf[0] = byte(v >> 8)
	buf[1] = byte(v)
	buf = append(buf, b.Bytes()...)
	_, err := conn.Write(buf)
	return err
}

func syncNodeInfo(conn net.Conn, mine *NodeInfo) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(mine); err != nil {
		return err
	}
	err := SendConnMsg(conn, &b)
	return err
}

func (n *Node) Dial(network, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return ext.Error(err)
	}

	if err := syncNodeInfo(conn, &n.NodeInfo); err != nil {
		return ext.Error(err)
	}
	return nil
}

func (n *Node) DialNode(target *Node) (err error) {
	return n.Dial(target.NodeInfo.Network, target.NodeInfo.String)
}

func (n *Node) accept() {
	for {
		conn, err := n.Listener.Accept()
		if err != nil {
			ext.Error(err)
			return
		}
		go n.dealOneCon(conn)
	}
}

func (n *Node) dealOneCon(conn net.Conn) {
	header := make([]byte, 2)
	var connName string
	isFirstData := false

	defer func() {
		conn.Close()
		n.delConn <- &connName
	}()

	for {
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(conn, header)
		if err != nil {
			ext.Error(err)
			return
		}

		size := binary.BigEndian.Uint16(header)
		data := make([]byte, size)
		//conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err = io.ReadFull(conn, data)
		if err != nil {
			ext.Error(err)
			return
		}

		b := bytes.NewBuffer(data)
		if !isFirstData {
			isFirstData = true
			dec := gob.NewDecoder(b)
			var r RemoteNode
			if err := dec.Decode(&(r.NodeInfo)); err != nil {
				ext.Error(err)
				continue
			}
			r.Conn = conn
			connName = r.NodeInfo.Name
			n.newConn <- &r
		} else {
			ok, msg := proto.DecodeMsg(b)
			if ok {
				_ = msg //消息处理
			}
		}
	}
}

func (n *Node) updateConnect() {
	for {
		select {
		case rc, ok := <-n.newConn:
			if !ok {

			}
			n.connected[rc.Name] = *rc
		case delName, ok := <-n.delConn:
			if !ok {

			}
			delete(n.connected, *delName)
		}
	}
}

func (n *Node) Connected() RemoteNodeMap {
	return n.connected
}
