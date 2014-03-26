package god

import (
	// "encoding"
	"proto"
)

type Marshaler interface {
	// encoding.BinaryMarshaler
	// encoding.BinaryUnmarshaler
}

type Handler interface {
	Handle(packID proto.PacketID, data Marshaler) (retID proto.PacketID, ret Marshaler, err error)
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
