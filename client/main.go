package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	service := "127.0.0.1:7777"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if checkError(err, "Resolve") {
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if checkError(err, "Dial") {
		os.Exit(3)
	}
	//启动客户端发送线程
	go chatSend(conn)

	//开始客户端轮训
	buf := make([]byte, 1024)
	for {
		len, err := conn.Read(buf)
		if checkError(err, "read") {
			if strings.EqualFold(err.Error(), "EOF") {
				continue
			}
		}
		fmt.Println(string(buf[0:len]))
	}
	defer conn.Close()
}

func checkError(err error, info string) (ret bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error[%s]: %s", info, err.Error())
		ret = true
	}
	ret = false
	return
}

////////////////////////////////////////////////////////
//
//客户端发送线程
//参数
//      发送连接 conn
//
////////////////////////////////////////////////////////
func chatSend(conn net.Conn) {

	var input string
	username := conn.LocalAddr().String()
	for {

		fmt.Scanln(&input)
		if input == "/quit" {
			fmt.Println("ByeBye..")
			conn.Close()
			os.Exit(0)
		}

		lens, err := conn.Write([]byte(username + " Say :::" + input))
		fmt.Println("Write Len: " + string(lens))
		if err != nil {
			fmt.Println("send " + err.Error())
			conn.Close()
			break
		}
	}
}

func handleConnect() {
}
