package god

import (
	"bytes"
	"ext"
	"net"
	"proto"
	"time"
)

type Remote struct {
	*process
	mes       Messenger
	conn      net.Conn
	nodeAddr  string
	objs      map[PID]string
	msgBuffer chan *proto.Message
}

type OperateObjs interface {
	AddRemoteObj(*proto.Message)
	RemoveRemoteObj(*proto.Message)
}

type SendObjsInfo struct {
	NodeAddr string
	Objs     []PID
}

type RemoteHandler struct {
	OperateObjs
}

func NewRemote(mes Messenger, conn net.Conn) Processor {
	r := new(Remote)
	r.process = NewProcess(m, new(RemoteHandler{r}), nil)
	r.mes = mes
	r.objs = make(map[PID]string, 100)
	r.msgBuffer = make(chan *proto.Message, CHAN_BUFF_NUM)
	err := r.dial(conn)
	if err != nil {
		return nil
	}
	go r.readRun()
	go r.writeRun()
	return r
}

func (r *Remote) dial(conn net.Conn) error {
	rObjs := syncNodeInfo(conn, GetNodeInfo(), r.AllProcessInfo())
	r.nodeAddr = rObjs.NodeAddr
	r.conn = conn
	for i := 0; i < len(rObjs.Objs); i++ {
		r.objs[rObjs.Objs[i]] = ""
	}
	return nil
}

func syncNodeInfo(conn net.Conn, nodeInfo NodeInfo, selfObjs []PID) *SendObjsInfo {
	defer nodeTrace.UT(nodeTrace.T("Node::syncNodeInfo\t%s\tto\t%s", nodeInfo.Name, conn.RemoteAddr().String()))

	var b, wb bytes.Buffer
	var err error

	enc := gob.NewEncoder(&b)
	ext.AssertE(enc.Encode(nodeInfo))
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

	var remote NodeInfo

	rb := bytes.NewBuffer(data)
	dec := gob.NewDecoder(rb)
	ext.AssertE(dec.Decode(&remote))

	var retObjs SendObjsInfo
	retObjs.NodeAddr = remote.String
	ext.AssertE(dec.Decode(&(retObjs.Objs)))

	return &retObjs
}

func (r *Remote) readRun() {
	defer r.conn.Close()
	header := make([]byte, 2)

	for {
		r.conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(r.conn, header)
		ext.AssertE(err)

		data := make([]byte, BYTE_ORDER.Uint16(header))
		_, err = io.ReadFull(r.conn, data)
		ext.AssertE(err)

		b := bytes.NewBuffer(data)
		ok, msg := proto.DecodeMsg(b)
		if ok {
			r.mes.Notify(msg.Reciever, &msg)
		}
	}
}

func (r *Remote) writeRun() {
	var buff bytes.Buffer
	for {
		select {
		case msg <- r.msgBuffer:
			buff.Reset()
			ret := proto.EncodeMsg(&buff, msg)
			if ret == false {
				ext.Errorf("Error enc Msg")
			} else {
				n, err := r.conn.Write(buff.Bytes())
				if err != nil {
					ext.Errorf("Error send bytes:", n, "Reason", err.Error())
				}
			}
		}
	}

}

func (r *Remote) proNotify(msg *proto.Message) error {
	if msg.Reciever == r.pid() {
		return r.processor.proNotify(packID, data)
	} else {
		r.msgBuffer <- msg
		return nil
	}
}

func (r *Remote) proCall(msg *proto.Message) (error, *proto.Message) {
	if msg.Reciever == r.pid() {
		return r.processor.proCall(packID, data)
	} else {
		r.msgBuffer <- msg
		return nil, nil
	}
}

func (r *Remote) AddRemoteObj(msg *proto.Message) {

}

func (r *Remote) RemoveRemoteObj(msg *proto.Message) {

}

func (r *RemoteHandler) Handle(packID proto.PacketID, data *proto.Message) {
	//TODO remote process add or remove
	//r.AddRemoteObj(data)
	//r.RemoveRemoteObj(data)
	return nil
}
