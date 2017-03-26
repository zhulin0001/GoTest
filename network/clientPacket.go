package network

import (
	"bytes"
	"encoding/binary"
)

//PacketClientHeader  MainCmd, SubCmd, BodyLen, Encrypt
type PacketClientHeader struct {
	MainCmd uint16
	SubCmd  uint16
	BodyLen uint32
	Encrypt uint8
}

//PacketClientHeaderSize = MainCmd + SubCmd + Length + Encrypt
func PacketClientHeaderSize() uint32 {
	return 2 + 2 + 4 + 1
}

//NewPacketClientHeader construc header with bytes
func NewPacketClientHeader(bytes []byte) *PacketClientHeader {
	if uint32(len(bytes)) < PacketClientHeaderSize() {
		return nil
	}
	var p = new(PacketClientHeader)
	p.MainCmd = binary.BigEndian.Uint16(bytes[:2])
	p.SubCmd = binary.BigEndian.Uint16(bytes[2:4])
	p.BodyLen = binary.BigEndian.Uint32(bytes[4:8])
	p.Encrypt = bytes[8]
	return p
}

//Bytes return PacketClientHeader in bytes
func (p *PacketClientHeader) Bytes() []byte {
	buf := make([]byte, PacketClientHeaderSize())
	binary.BigEndian.PutUint16(buf[:2], p.MainCmd)
	binary.BigEndian.PutUint16(buf[2:4], p.SubCmd)
	binary.BigEndian.PutUint32(buf[4:8], p.BodyLen)
	buf[8] = p.Encrypt
	return buf
}

//PacketClient use for communation
type PacketClient struct {
	Header *PacketClientHeader
	Body   []byte
}

//Bytes return PacketClient in bytes
func (p *PacketClient) Bytes() []byte {
	buf := make([][]byte, PacketClientHeaderSize()+p.Header.BodyLen)
	buf[0] = p.Header.Bytes()
	buf[1] = p.Body
	return bytes.Join(buf, []byte(""))
}
