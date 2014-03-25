package god

import (
// "encoding"
)

type Marshaler interface {
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler
}

type Handler interface {
	Handle(packID PacketID, data Marshaler) (retID PacketID, ret Marshaler, err error)
}

type PacketID int64

type Message struct {
	Sender UUID
	Data   interface{}
	PackID PacketID
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
