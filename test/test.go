package main

import (
	"common"
	"encoding/binary"
	"fmt"
)

func main() {
	headerBuf := make([]byte, common.NetPacketHeaderSize())
	headerBuf[0] = 'a'
	header := common.NewNetPacketHeader(headerBuf)

	body := []byte("hello")

	fmt.Printf("hello %X\n", body)

	p := &common.NetPacket{Header: header, Body: body}
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
