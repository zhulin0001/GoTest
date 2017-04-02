package main

import (
	"config"
	"fmt"
	"io/ioutil"
	"net"
	"network"
	"os"
	"utils"

	"strconv"

	"runtime"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// net.Conn.RemoteAddr() => true/false
var connMap = make(map[string]bool)

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
	cMsgQueue := make(chan *network.PacketWrapper, channelBufNum)
	cConnList := make([]*network.ConnHandler, conf.Server.MaxClient)
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

func dispatchMsg(mq chan *network.PacketWrapper) {
	for {
		select {
		case msg := <-mq:
			conn := msg.RawCon
			packet := msg.Packet
			remoteAddr := conn.RemoteAddr().String()
			if v, ok := connMap[remoteAddr]; ok {
				//未认证过的
				if !v {
					packet.Header.MainCmd
				}
			}
			log.Info("Dispatch: ", remoteAddr, packet.Bytes())
		}
	}
}

func startListen(addr string, rc chan *network.PacketWrapper, hl []*network.ConnHandler) {
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
		remoteAddr := conn.RemoteAddr().String()

		if _, ok := connMap[remoteAddr]; ok {
			log.Info(remoteAddr, " Reconnect")
		}
		//所有新监听的连接都置为不受信任
		connMap[remoteAddr] = false

		log.Info("connection is connected from ...", conn.RemoteAddr())
		connHandle := new(network.ConnHandler)
		connHandle.RawCon = conn
		connHandle.ErrChan = make(chan *network.InternalChannelMsg)
		connHandle.ReadChan = rc
		connHandle.WriteChan = make(chan *network.PacketClient)
		connHandle.CloseChan = make(chan bool)
		hl = append(hl, connHandle)

		go connReadLoop(connHandle)
		go connWriteLoop(connHandle)
	}
}

func connReadLoop(handler *network.ConnHandler) {
	var conn = handler.RawCon
	for {
		select {
		case errChan := <-handler.ErrChan:
			log.Warn("Recive Error: ", errChan.String())
			handler.CloseChan <- true
		default:
			pWrapper, err := handler.Parser.Decode(conn)
			if utils.CheckError(err, "Parse Packet") {
				handler.CloseChan <- true
			}
			handler.ReadChan <- pWrapper
		}
	}
}

func connWriteLoop(handler *network.ConnHandler) {
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
