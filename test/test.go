package main

import (
	"common"
	"encoding/binary"
	"fmt"
)

func main() {
	header := new(common.NetPacketHeader)
	header.MainCmd = 106
	header.SubCmd = 2
	header.BodyLen = 3
	header.Encrypt = 1

	body := []byte("hello")

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
