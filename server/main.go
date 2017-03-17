package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := "127.0.0.1:7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if checkError(err) {
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if checkError(err) {
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handlerClient(conn)
	}
}

func checkError(err error) (ret bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
		ret = true
	}
	ret = false
	return
}

func handlerClient(conn net.Conn) {
	fmt.Fprintf(os.Stderr, "New Connection From: %s", conn.RemoteAddr().String())
	for {
		daytime := time.Now().String()
		conn.Write([]byte(daytime))
		buf := make([]byte, 1024)
		readLen, err := conn.Read(buf)
		if checkError(err) {
			conn.Close()
			break
		}
		if readLen > 0 {
			fmt.Fprintf(os.Stderr, "Recv: %s", string(buf[0:readLen]))
		} else {
			conn.Close()
			break
		}
	}
}
