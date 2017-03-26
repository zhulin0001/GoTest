package network

import (
	"fmt"
	"net"
)

//PacketWrapper for internal use
type PacketWrapper struct {
	RawCon net.Conn
	Packet *PacketClient
}

//PacketParser 包解析器
type PacketParser interface {
	Decode(net.Conn) (*PacketWrapper, error)
}

// ClientPacketParser defines a special parser for client
type ClientPacketParser struct{}

//Decode for ClientPacketParser
func (parser ClientPacketParser) Decode(conn net.Conn) (p *PacketWrapper, err error) {
	headerBuf := make([]byte, PacketClientHeaderSize())
	for {
		_, err = conn.Read(headerBuf)
		if err != nil {
			err = fmt.Errorf("Read Header Error: %s", err.Error())
			break
		}
		header := NewPacketClientHeader(headerBuf)
		if header == nil {
			err = fmt.Errorf("Read Invalid Header: %X", headerBuf)
			break
		}
		bodyData := make([]byte, header.BodyLen) //考虑长度为0
		if header.BodyLen > 0 {
			_, err = conn.Read(bodyData)
			if err != nil {
				err = fmt.Errorf("Read PacketBody Error: %s", err.Error())
				break
			}
		}
		packet := &PacketClient{Header: header, Body: bodyData}
		p = &PacketWrapper{RawCon: conn, Packet: packet}
		break
	}
	return
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
	WriteChan chan *PacketClient
	CloseChan chan bool
	Parser    PacketParser
}

func (h *ConnHandler) start(conn net.Conn) {

}
