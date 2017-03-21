package common

//NetPacket use for communation
type NetPacket struct{
	MainCmd	uint16
	SubCmd	uint16
	Length	uint32
	Encrypt	uint8
	Body	[]byte
}

//InternalChannelMsg use for Service know what error occurs and stop self
type InternalChannelMsg struct {
	Code		int32
	Reason		string
	MoreInfo 	interface{}
}