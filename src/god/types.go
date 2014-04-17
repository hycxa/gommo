package god

import (
	// "encoding"
	"proto"
)

const (
	CHAN_BUFF_NUM = 16 //chan buffer max deal num
	TCP_TIMEOUT   = 60 // tcp read timeout
)

type Marshaler interface {
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler
}

type NodeInfo struct {
	Name     string
	Network  string
	String   string
	NodeType string
}

type Handler interface {
	Handle(packID proto.PacketID, data *proto.Message) error
}

type NotifyFun interface {
	notify(data *proto.Message)
}

// func (m Marshaler) MarshalBinary() ([]byte, error) {
// 	var b bytes.Buffer
// 	enc := gob.NewEncoder(&b)
// 	err := enc.Encode(m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return b.Bytes(), nil
// }

// func (m Marshaler) UnmarshalBinary(data []byte) error {
// 	b := bytes.NewBuffer(data)
// 	dec := gob.NewDecoder(b)
// 	return dec.Decode(m)
// }
