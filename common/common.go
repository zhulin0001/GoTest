package common

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

//NetPacketHeader: MainCmd, SubCmd, BodyLen, Encrypt
type NetPacketHeader struct {
	MainCmd uint16
	SubCmd  uint16
	BodyLen uint32
	Encrypt uint8
}

//NetPacketHeaderSize = MainCmd + SubCmd + Length + Encrypt
func NetPacketHeaderSize() uint32 {
	return 2 + 2 + 4 + 1
}

//NewNetPacketHeader construc header with bytes
func NewNetPacketHeader(bytes []byte) *NetPacketHeader {
	if uint32(len(bytes)) < NetPacketHeaderSize() {
		return nil
	}
	var p = new(NetPacketHeader)
	p.MainCmd = binary.BigEndian.Uint16(bytes[:2])
	p.SubCmd = binary.BigEndian.Uint16(bytes[2:4])
	p.BodyLen = binary.BigEndian.Uint32(bytes[4:8])
	p.Encrypt = bytes[8]
	return p
}

//Bytes return NetPacketHeader in bytes
func (p *NetPacketHeader) Bytes() []byte {
	buf := make([]byte, NetPacketHeaderSize())
	binary.BigEndian.PutUint16(buf[:2], p.MainCmd)
	binary.BigEndian.PutUint16(buf[2:4], p.SubCmd)
	binary.BigEndian.PutUint32(buf[4:8], p.BodyLen)
	buf[8] = p.Encrypt
	return buf
}

//NetPacket use for communation
type NetPacket struct {
	Header *NetPacketHeader
	Body   []byte
}

//Bytes return NetPacket in bytes
func (p *NetPacket) Bytes() []byte {
	buf := make([][]byte, NetPacketHeaderSize()+p.Header.BodyLen)
	buf[0] = p.Header.Bytes()
	buf[1] = p.Body
	return bytes.Join(buf, []byte(""))
}

//PacketWrapper for internal use
type PacketWrapper struct {
	RawCon net.Conn
	Packet *NetPacket
}

//InternalChannelMsg use for Service know what error occurs and stop self
type InternalChannelMsg struct {
	Code     int32
	Reason   string
	MoreInfo interface{}
}

//String use for log the detail
func (p *InternalChannelMsg) String() string {
	return fmt.Sprintf("Code:[%d], Reason:\"%s\"", p.Code, p.Reason)
}

//ConnHandler use for manage connections
type ConnHandler struct {
	RawCon    net.Conn
	ErrChan   chan *InternalChannelMsg
	ReadChan  chan *PacketWrapper
	WriteChan chan *NetPacket
	CloseChan chan bool
}

func (h *ConnHandler) start(conn net.Conn) {

}
