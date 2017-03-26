package network

//PacketInternalHeader 内部通讯包的包头格式
type PacketInternalHeader struct {
	ServerID uint16
	Type     uint8
}

//PacketInternal for server communation
type PacketInternal struct {
	Header *PacketInternalHeader
	Body   []byte
}
