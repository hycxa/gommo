package god

import (
	// "encoding"
	"proto"
)

const (
	CHAN_BUFF_NUM    = 16 //chan buffer max deal num
	TCP_TIMEOUT      = 60 // tcp read timeout
	WORKER_NUM_LIMIT = 4  //worker num limit
)

type Marshaler interface {
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler
}

type NodeInfo struct {
	Name    string
	Network string
	String  string
}

type Handler interface {
	Handle(msg *proto.Message) error
}

type WorkerNotifyFun interface {
	notify(msg *proto.Message)
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
