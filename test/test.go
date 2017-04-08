package main

import (
	"encoding/binary"
	"fmt"
	"network"
)

func main1() {
	headerBuf := make([]byte, network.PacketClientHeaderSize())
	headerBuf[0] = 'a'
	header := network.NewPacketClientHeader(headerBuf)

	body := []byte("hello")

	fmt.Printf("hello %X\n", body)

	p := &network.PacketClient{Header: header, Body: body}
	fmt.Println(header.Bytes())
	fmt.Printf("%x\n", p.Bytes())

	number := 2
	fmt.Println(&number)
	test := make([]byte, 2)
	binary.BigEndian.PutUint16(test, 2)
	fmt.Println(test)
	binary.LittleEndian.PutUint16(test, 2)
	fmt.Println(test)
}

func main() {
	Depushe()
}
