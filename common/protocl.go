package common

import (
	"bytes"
	"encoding/binary"
)

const (
	ConstHeader         = "{28285252_hcaiyue_top_15159898}"
	ConstHeaderLength   = len(ConstHeader)
	ConstSaveDataLength = 8
	ConstFooter         = "{137621_woyaoyigebaojieshula_159687}"
	ConstFooterLength   = len(ConstFooter)
)

//封包
func Packet(message []byte) []byte {
	var buf bytes.Buffer
	buf.Write([]byte(ConstHeader))
	buf.Write(IntToBytes(len(message)))
	buf.Write(message)
	buf.Write([]byte(ConstFooter))
	return buf.Bytes()
}

//解包
func Unpack(buffer []byte, reader *[]byte) []byte {
	length := len(buffer)
	biao := 0
	if length <= ConstHeaderLength+ConstSaveDataLength+ConstFooterLength {
		return buffer
	}

	if !bytes.Contains(buffer, []byte(ConstHeader)) || !bytes.Contains(buffer, []byte(ConstFooter)) {
		return buffer
	}
	start := bytes.Index(buffer, []byte(ConstHeader))
	if start == -1 {
		return buffer
	}
	end := bytes.Index(buffer[start+ConstHeaderLength:], []byte(ConstFooter))
	if end == -1 {
		return buffer
	}
	end += start + ConstHeaderLength
	biao = start + ConstHeaderLength
	messageLength := BytesToInt(buffer[biao : biao+ConstSaveDataLength])
	if end-start-ConstHeaderLength-ConstSaveDataLength != messageLength {
		return buffer[end+ConstFooterLength:]
	}
	biao += ConstSaveDataLength
	*reader = buffer[biao : biao+messageLength]
	biao += messageLength + ConstFooterLength
	if biao == length {
		return make([]byte, 0)
	}
	return buffer[biao:]
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int64(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
