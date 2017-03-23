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

	"runtime"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	conf := readConfig()
	if conf == nil {
		log.Error("Read Config Failed!")
		os.Exit(1)
	}
	//创建客户端读写channe
	channelBufNum := conf.Server.MaxChannelBuf
	cMsgQueue := make(chan *common.PacketWrapper, channelBufNum)
	cConnList := make([]*common.ConnHandler, conf.Server.MaxClient)
	fmt.Printf("Cap[%d], Len[%d]\n", cap(cConnList), len(cConnList))

	go dispatchMsg(cMsgQueue)
	serverAddrForClients := conf.Server.ListenIP + ":" + strconv.Itoa(conf.Server.Port)
	go startListen(serverAddrForClients, cMsgQueue, cConnList)

	// sMsgQueue := make(chan []byte, 10)
	// swc := make(chan []byte, 10)
	// sConnList := make([]common.ConnHandler, conf.Server.MaxClient)
	// serverAddrForServers := conf.Server.ListenIP + ":" + strconv.Itoa(conf.Server.PortInternal)
	// go startListen(serverAddrForServers, sMsgQueue, &sConnList)

	for {
		runtime.Gosched()
	}
}

func readConfig() (cfg *config.ACConfig) {
	var configFileName = "config.toml"
	if len(os.Args) > 1 {
		configFileName = os.Args[1]
	}
	var confData []byte
	if utils.CheckFileIsExist(configFileName) {
		var err error
		confData, err = ioutil.ReadFile(configFileName)
		if utils.CheckError(err, "ReadFile") {
			log.Info("Read File[%s] Error: ", configFileName, err.Error())
		}
	} else {
		fmt.Println("文件不存在")
	}
	var conf = new(config.ACConfig)
	err := toml.Unmarshal(confData, conf)
	if utils.CheckError(err, "toml.Decode") {
		log.Info("toml decode error: " + err.Error())
	} else {
		cfg = conf
	}
	return cfg
}

func dispatchMsg(mq chan *common.PacketWrapper) {
	for {
		select {
		case msg := <-mq:
			conn := msg.RawCon
			packet := msg.Packet
			log.Info("Dispatch: ", conn.RemoteAddr(), packet.Bytes())
		}
	}
}

func startListen(addr string, rc chan *common.PacketWrapper, hl []*common.ConnHandler) {
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
			break
		}
		log.Info("connection is connected from ...", conn.RemoteAddr())
		connHandle := new(common.ConnHandler)
		connHandle.RawCon = conn
		connHandle.ErrChan = make(chan *common.InternalChannelMsg)
		connHandle.ReadChan = rc
		connHandle.WriteChan = make(chan *common.NetPacket)
		connHandle.CloseChan = make(chan bool)
		hl = append(hl, connHandle)

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
			bodyData := make([]byte, header.BodyLen) //考虑长度为0
			if header.BodyLen > 0 {
				_, err = conn.Read(bodyData)
				if utils.CheckError(err, "ReadLoop") {
					log.Info("Read PacketBody Error: ", err.Error())
					handler.CloseChan <- true
					break
				}
			}
			packet := &common.NetPacket{Header: header, Body: bodyData}
			pWrapper := &common.PacketWrapper{RawCon: handler.RawCon, Packet: packet}
			handler.ReadChan <- pWrapper
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
