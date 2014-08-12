package god

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"ext"
	"io"
)

func WriteBytes(w io.Writer, data []byte) {
	ext.AssertE(binary.Write(w, BYTE_ORDER, uint16(len(data))))
	ext.AssertE(binary.Write(w, BYTE_ORDER, data))
	ext.LogDebug("WRITTEN\t%d", len(data))
}

func ReadBytes(r io.Reader) []byte {
	var size uint16
	ext.AssertE(binary.Read(r, BYTE_ORDER, &size))
	data := make([]byte, size)
	ext.AssertE(binary.Read(r, BYTE_ORDER, data))
	ext.LogDebug("READ\t%d", size)
	return data
}

func GobDecode(bin []byte, data interface{}) {
	buf := bytes.NewBuffer(bin)
	decoder := gob.NewDecoder(buf)
	ext.AssertE(decoder.Decode(data))
}

func GobEncode(data interface{}) []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	ext.AssertE(encoder.Encode(data))
	return buf.Bytes()
}
