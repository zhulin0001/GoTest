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
		os.Exit(3)
	}
	conns := make(map[string]chan string)
	messages := make(chan string, 10000)
	addrs := make(chan string, 10000)
	// //启动服务器广播线程
	// go echoHandler(&conns, messages)
	go readMsg(messages, addrs, &conns)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		msg := make(chan string)
		conns[conn.RemoteAddr().String()] = msg
		fmt.Println("connection is connected from ...", conn.RemoteAddr().String())
		go clientRead(conn, messages, addrs)
		go clientWrite(conn, msg, addrs)
	}
}

func checkError(err error) (ret bool) {
	ret = false
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
		ret = true
	}
	return
}

func clientRead(conn net.Conn, messages chan string, addrs chan string) {
	buf := make([]byte, 1024)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err) {
			addrs <- conn.RemoteAddr().String()
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}
		reciveStr := string(buf[0:lenght])
		fmt.Println("Rec[" + conn.RemoteAddr().String() + "] Say : " + reciveStr)
		messages <- reciveStr
	}
}

func clientWrite(conn net.Conn, messages chan string, addrs chan string) {
	for {
		msg := <-messages
		_, err := conn.Write([]byte(msg))
		if checkError(err) {
			addrs <- conn.RemoteAddr().String()
			conn.Close()
			break
		}
	}
}

func readMsg(messages chan string, addrs chan string, conns *map[string]chan string) {
	for {
		select {
		case str := <-messages:
			for _, value := range *conns {
				value <- str
			}
		case addr := <-addrs:
			delete(*conns, addr)
			fmt.Println("remove " + addr)
		}
	}
}

////////////////////////////////////////////////////////
//
//服务器发送数据的线程
//
//参数
//      连接字典 conns
//      数据通道 messages
//
////////////////////////////////////////////////////////
func echoHandler(conns *map[string]net.Conn, messages chan string) {

	for {
		msg := <-messages
		if len(msg) > 0 {
			fmt.Println(msg)
			for key, value := range *conns {
				fmt.Println("connection is connected from ...", key)
				_, err := value.Write([]byte(msg))
				if err != nil {
					fmt.Println(err.Error())
					delete(*conns, key)
					return
				}
			}
		}
	}
}

////////////////////////////////////////////////////////
//
//服务器端接收数据线程
//参数：
//      数据连接 conn
//      通讯通道 messages
//
////////////////////////////////////////////////////////
func Handler(conn net.Conn, messages chan string) {

	fmt.Println("connection is connected from ...", conn.RemoteAddr().String())

	buf := make([]byte, 1024)
	for {
		lenght, err := conn.Read(buf)
		if checkError(err) == true {
			conn.Close()
			break
		}
		if lenght > 0 {
			buf[lenght] = 0
		}
		fmt.Println("Rec[", conn.RemoteAddr().String(), "] Say :", string(buf[0:lenght]))
		reciveStr := string(buf[0:lenght])
		messages <- reciveStr
	}
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
