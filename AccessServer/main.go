package main

import (
	"common"
	"config"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"utils"

	"strconv"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	var configFileName = "config.toml"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	var conf = new(config.ACConfig)
	var confData []byte
	if utils.CheckFileIsExist(configFileName) {
		var err error
		confData, err = ioutil.ReadFile(configFileName)
		if utils.CheckError(err, "ReadFile") {
			log.Info("Read File[%s] Error: ", configFileName, err.Error())
		}
		//fmt.Println("文件存在")
		//fmt.Println(confData)
	} else {
		fmt.Println("文件不存在")
	}
	err := toml.Unmarshal(confData, conf)
	if utils.CheckError(err, "toml.Decode") {
		log.Info("toml decode error: " + err.Error())
		os.Exit(0)
	}
	//创建客户端读写channe
	channelBufNum := conf.Server.MaxChannelBuf
	crc := make(chan []byte, channelBufNum)
	cwc := make(chan []byte, channelBufNum)

	serverAddrForClients := conf.Server.ListenIP + ":" + strconv.Itoa(conf.Server.Port)
	go startListenForClients(serverAddrForClients, crc, cwc)

	src := make(chan []byte, 10)
	swc := make(chan []byte, 10)
	serverAddrForServers := conf.Server.ListenIP + ":" + strconv.Itoa(conf.Server.PortInternal)
	startListenForServers(serverAddrForServers, src, swc)
}

func startListenForClients(addr string, rc chan []byte, wc chan []byte) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if utils.CheckError(err, "Resolve") {
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	defer listener.Close()
	if utils.CheckError(err, "Listen") {
		os.Exit(1)
	}
	log.Info("Server Listening On ", tcpAddr)

	for {
		conn, err := listener.Accept()
		if utils.CheckError(err, "OnListen") {
			log.Error("Socket Error On Accept: ", err.Error())
		}
		log.Info("connection is connected from ...", conn.RemoteAddr().String())
	}
}

func startListenForServers(addr string, rc chan []byte, wc chan []byte) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if utils.CheckError(err, "Resolve") {
		os.Exit(1)
	}
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	defer listener.Close()
	if utils.CheckError(err, "Listen") {
		os.Exit(1)
	}
	log.Info("Server Listening On ", tcpAddr)

	for {
		conn, err := listener.Accept()
		if utils.CheckError(err, "OnListen") {
			log.Error("Socket Error On Accept: ", err.Error())
		}
		log.Info("connection is connected from ...", conn.RemoteAddr().String())
		connHandle := new(common.ConnHandler)
		connHandle.RawCon = conn
		connHandle.ErrChan = make(chan *common.InternalChannelMsg)
		connHandle.ReadChan = make(chan *common.NetPacket)
		connHandle.WriteChan = make(chan *common.NetPacket)
		connHandle.CloseChan = make(chan bool)

		go connReadLoop(connHandle)
		go connWriteLoop(connHandle)
	}
}

func connReadLoop(handler *common.ConnHandler) {
	var conn = handler.RawCon
	for {
		select {
		case errChan := <-handler.ErrChan:
			log.Warn("Recive Error: ", errChan.String())
			handler.CloseChan <- true
		default:
			headerBuf := make([]byte, common.NetPacketHeaderSize())
			_, err := conn.Read(headerBuf)
			if utils.CheckError(err, "ReadLoop") {
				log.Info("Read Header Error: ", err.Error())
				handler.CloseChan <- true
				break
			}
			header := common.NewNetPacketHeader(headerBuf)
			if header == nil {
				log.Warn("Read Invalid Header: ", headerBuf)
				handler.CloseChan <- true
				break
			}
			bodyData := make([]byte, header.BodyLen)
			_, err = conn.Read(bodyData)
			if utils.CheckError(err, "ReadLoop") {
				log.Info("Read PacketBody Error: ", err.Error())
				handler.CloseChan <- true
				break
			}
			packet := &common.NetPacket{Header: header, Body: bodyData}
			handler.ReadChan <- packet
		}
	}
}

func connWriteLoop(handler *common.ConnHandler) {
	var conn = handler.RawCon
	for {
		select {
		case errChan := <-handler.ErrChan:
			log.Warn("Recive Error: ", errChan.String())
			handler.CloseChan <- true
		case packet := <-handler.WriteChan:
			_, err := conn.Write(packet.Bytes())
			if utils.CheckError(err, "WriteLoop") {
				log.Info("Write Packet Error: ", err.Error())
				handler.CloseChan <- true
				break
			}
		}
	}
}
