package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	go handleConnect()
	service := "127.0.0.1:7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	server, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	for i := 0; i < 10; i++ {
		buf := make([]byte, 1024)
		readLen, err := server.Read(buf)
		checkError(err)
		if readLen > 0 {
			fmt.Fprintf(os.Stderr, "Recv: %s \n", string(buf[0:readLen]))
		}
		server.Write([]byte("Hello world"))
	}
	defer server.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
		os.Exit(1)
	}
}

func handleConnect() {
	service := "127.0.0.1:7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	server, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	for i := 0; i < 10; i++ {
		buf := make([]byte, 1024)
		readLen, err := server.Read(buf)
		checkError(err)
		if readLen > 0 {
			fmt.Fprintf(os.Stderr, "Recv: %s \n", string(buf[0:readLen]))
		}
		server.Write([]byte("Hello world"))
	}
	defer server.Close()
}
