package god

import (
	"bytes"
	"ext"
	"io"
	"net"
	"proto"
	"time"
)

type Agent struct {
	*process
	mes         Messenger
	nFun        WorkerNotifyFun
	conn        net.Conn
	writeBuffer chan *proto.Message
}

func NewAgent(mes Messenger, nFun WorkerNotifyFun, conn net.Conn) Processor {
	a := new(Agent)
	a.process = NewProcess(mes, a)
	a.mes = mes
	a.nFun = nFun
	a.conn = conn
	a.writeBuffer = make(chan *proto.Message, CHAN_BUFF_NUM)
	go a.readRun()
	go a.writeRun()
	return a
}

func (a *Agent) readRun() {
	defer a.conn.Close()
	header := make([]byte, 2)

	for {
		a.conn.SetReadDeadline(time.Now().Add(TCP_TIMEOUT * time.Second))
		_, err := io.ReadFull(a.conn, header)
		ext.AssertE(err)

		data := make([]byte, BYTE_ORDER.Uint16(header))
		_, err = io.ReadFull(a.conn, data)
		ext.AssertE(err)

		b := bytes.NewBuffer(data)
		ok, msg := proto.DecodeMsg(b)
		if ok {
			a.nFun.notify(msg)
		}
	}
}

func (a *Agent) Handle(msg *proto.Message) error {
	msgType := proto.GetPacketScope(msg.PacketID)
	if msgType == proto.PACKAGE_SYSTEM {
		//control operate
		return nil
	} else if msgType == proto.PACKAGE_USER {
		a.nFun.notify(msg)
		return nil
	} else {
		return ext.MyError{"unknow msgType"}
	}
}

//此接口给lua用的
func (a *Agent) writeMsg(msg *proto.Message) {
	a.writeBuffer <- msg
}

func (r *Agent) writeRun() {
	var buff bytes.Buffer
	for {
		select {
		case msg := <-r.writeBuffer:
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
