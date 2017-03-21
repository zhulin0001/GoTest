package main

import (
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
			log.Debug("Read File[%s] Error: ", configFileName, err.Error())
		}
		//fmt.Println("文件存在")
		//fmt.Println(confData)
	} else {
		fmt.Println("文件不存在")
	}
	err := toml.Unmarshal(confData, conf)
	if utils.CheckError(err, "toml.Decode") {
		log.Debug("toml decode error: " + err.Error())
		os.Exit(0)
	}
	//创建客户端读写channe
	channelBufNum := conf.Server.MaxChannelBuf
	crc := make(chan []byte, channelBufNum)
	cwc := make(chan []byte, channelBufNum)

	serverAddrForClients := conf.Server.ListenIP + ":" + strconv.Itoa(conf.Server.Port)
	startListenForClients(serverAddrForClients, crc, cwc)

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

	for {
		conn, err := listener.Accept()
		if utils.CheckError(err, "OnListen") {
			os.Exit(3)
		}
		log.Debug("connection is connected from ...", conn.RemoteAddr().String())
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

	for {
		conn, err := listener.Accept()
		if utils.CheckError(err, "OnListen") {
			os.Exit(3)
		}
		log.Debug("connection is connected from ...", conn.RemoteAddr().String())
	}
}
