package god

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"ext"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"proto"
)

const (
	TCP_TIMEOUT = 12000
)

const (
	NODE_GS_TYPE   = "GS"
	NODE_GATE_TYPE = "GATE"
)

type NodeID hash.Hash

type NodeInfo struct {
	Name     string
	Network  string
	String   string
	NodeType string
}

type RemoteNode struct {
	NodeInfo
	net.Conn
}

type RemoteNodeMap map[string]*RemoteNode

type SendObjsInfo struct {
	NodeAddr string
	Objs     []proto.UUID
}

type Node struct {
	NodeInfo
	net.Listener
	newRemote       chan *RemoteNode
	newRemoteObjs   chan *SendObjsInfo
	closeRemoteObjs chan *[]proto.UUID
	closingNodes    chan string
	closeRequest    chan bool
	connected       RemoteNodeMap
	objects         map[proto.UUID]*Process
	remoteObjs      map[proto.UUID]string
}

var (
	nodeTrace = ext.Trace(false)
)

func NewNode(name, network, address string, nodeType string) *Node {
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
	self.NodeType = nodeType
	self.newRemote = make(chan *RemoteNode, 16)
	self.closingNodes = make(chan string, 16)
	self.newRemoteObjs = make(chan *SendObjsInfo)
	self.closeRemoteObjs = make(chan *[]proto.UUID)
	self.connected = make(map[string]*RemoteNode)
	self.closeRequest = make(chan bool)
	self.objects = make(map[proto.UUID]*Process)
	self.remoteObjs = make(map[proto.UUID]string)
	go self.accept()
	go self.update()
	return self
}

func syncNodeInfo(conn net.Conn, nodeInfo NodeInfo, objs *map[proto.UUID]*Process) (*RemoteNode, *SendObjsInfo) {
	defer nodeTrace.UT(nodeTrace.T("Node::syncNodeInfo\t%s\tto\t%s", nodeInfo.Name, conn.RemoteAddr().String()))

	var b, wb bytes.Buffer
	var err error

	enc := gob.NewEncoder(&b)
	ext.AssertE(enc.Encode(nodeInfo))

	selfObjs := make([]proto.UUID, len(*objs))
	var index = 0
	for uuid, _ := range *objs {
		selfObjs[index] = uuid
		index++
	}
	ext.AssertE(enc.Encode(selfObjs))

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

	var retObjs SendObjsInfo
	retObjs.NodeAddr = r.String
	ext.AssertE(dec.Decode(&(retObjs.Objs)))

	return &r, &retObjs
}

func (self *Node) Dial(network, address string) error {
	_, ok := self.connected[address]
	if ok || address == self.String {
		return nil
	}
	conn, err := net.Dial(network, address)
	if err != nil {
		return ext.LogError(err)
	}

	r, rObjs := syncNodeInfo(conn, self.NodeInfo, &self.objects)
	self.connected[r.String] = r
	for i := 0; i < len(rObjs.Objs); i++ {
		self.remoteObjs[rObjs.Objs[i]] = r.String
	}

	return nil
}

func (self *Node) DialNode(target *Node) (err error) {
	return self.Dial(target.NodeInfo.Network, target.NodeInfo.String)
}

func (self *Node) ConnOtherSvr() error {
	client := &http.Client{}
	reqPost, err := http.NewRequest("POST", "http://127.0.0.1:20000/locatePost", nil)
	if err != nil {
		return ext.LogError(err)
	}
	reqPost.Header.Set("Node-Addr", self.String)
	reqPost.Header.Set("service", self.NodeType)
	postRep, err := client.Do(reqPost)
	defer postRep.Body.Close()
	if err != nil {
		return ext.LogError(err)
	}

	reqGet, err := http.NewRequest("GET", "http://127.0.0.1:20000/locateGet", nil)
	if err != nil {
		return ext.LogError(err)
	}
	getSvrTab := make([]string, 2)
	getSvrTab[0] = NODE_GS_TYPE
	if self.NodeType == NODE_GS_TYPE {
		getSvrTab[0] = NODE_GS_TYPE
		getSvrTab[1] = NODE_GATE_TYPE
	}
	for i := 0; i < len(getSvrTab); i++ {
		if getSvrTab[i] != "" {
			reqGet.Header.Set("service", getSvrTab[i])
			resp, err := client.Do(reqGet)
			defer resp.Body.Close()
			if err != nil {
				return ext.LogError(err)
			}

			if resp.StatusCode == http.StatusOK {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return ext.LogError(err)
				}
				var svrList []string
				err = json.Unmarshal(body, &svrList)
				if err != nil {
					return ext.LogError(err)
				}
				for index := 0; index < len(svrList); index++ {
					err := self.Dial("tcp", svrList[index])
					if err != nil {
						ext.LogError(err)
					}
				}
			}

		}
	}
	return err
}

func (self *Node) accept() {
	for {
		conn, err := self.Accept()
		if err != nil {
			ext.LogError(err)
			return
		}
		r, rObjs := syncNodeInfo(conn, self.NodeInfo, &self.objects)
		self.newRemote <- r
		self.newRemoteObjs <- rObjs
		go self.dealOneCon(conn, r.String, self.closingNodes)
	}
}

func (self *Node) dealOneCon(conn net.Conn, nodeAddr string, closingNodes chan string) {
	header := make([]byte, 2)

	defer func() {
		conn.Close()
		closingNodes <- nodeAddr
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

func (self *Node) update() {
	for {
		select {
		case newRemote, ok := <-self.newRemote:
			if ok {
				self.connected[newRemote.String] = newRemote
			}
		case newObjs, ok := <-self.newRemoteObjs:
			if ok {
				for i := 0; i < len(newObjs.Objs); i++ {
					self.remoteObjs[newObjs.Objs[i]] = newObjs.NodeAddr
				}
			}
		case nodeAddr, ok := <-self.closingNodes:
			if ok {
				delete(self.connected, nodeAddr)
				var delTab []proto.UUID
				for uuid, addr := range self.remoteObjs {
					if addr == nodeAddr {
						delTab = append(delTab, uuid)
					}
				}
				for i := 0; i < len(delTab); i++ {
					delete(self.remoteObjs, delTab[i])
				}
			}
		case closeObjs, ok := <-self.closeRemoteObjs:
			if ok {
				for i := 0; i < len(*closeObjs); i++ {
					delete(self.remoteObjs, (*closeObjs)[i])
				}
			}
		case <-self.closeRequest:
			for _, remoteNode := range self.connected {
				remoteNode.Close()
			}
			return
		}
	}
}

func (self *Node) Connected() RemoteNodeMap {
	return self.connected
}

func (self *Node) Close() {
	self.closeRequest <- true
}

func (self *Node) Notify(source proto.UUID, target proto.UUID, packetID proto.PacketID, data interface{}) error {
	defer ext.UT(ext.T("NOTIFY"))
	t, ok1 := self.objects[target]
	otherSvrAddr, ok2 := self.remoteObjs[target]
	if !ok1 && !ok2 {
		return fmt.Errorf("Target %v is not found!", target)
	}
	m := proto.Message{Sender: source, Data: data, PackID: packetID}

	if ok1 {
		t.mq <- m
	} else {
		otherSvr, ok3 := self.connected[otherSvrAddr]
		if !ok3 {
			return fmt.Errorf("Target %v is not found!", target)
		}
		var b, wb bytes.Buffer
		ret := proto.EncodeMsg(&b, &m)
		if !ret {
			return fmt.Errorf("Encode fail %v", target)
		}
		ext.AssertE(binary.Write(&wb, BYTE_ORDER, uint16(len(b.Bytes()))))
		_, err := otherSvr.Conn.Write(wb.Bytes())
		ext.AssertE(err)
		_, err = otherSvr.Conn.Write(b.Bytes())
		ext.AssertE(err)
	}
	return nil
}

func (self *Node) Call(packID proto.PacketID, data proto.Message) (retID proto.PacketID, ret proto.Message, err error) {
	return 0, proto.Message{}, nil
}

func (self *Node) AddProcess(o *Process) {
	self.objects[o.UUID] = o
}
