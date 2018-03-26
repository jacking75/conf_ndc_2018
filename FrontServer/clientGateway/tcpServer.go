package clientGateway

import (
	"net"
	"os"

	"../monitoring"
	def "../typeDefine"
	"../utils"

	"github.com/emirpasic/gods/lists/arraylist"
)

var mClientListener *net.TCPListener

func TCPServerStart(config *def.Config) {
	defer monitoring.ServerStatusDestoryChannel()
	defer utils.PrintPanicStack()

	backEndList := getbackendServerList()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.ClientAddress)
	checkError(err)

	mClientListener, err = net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	utils.Logger.Info("listening on:", config.ClientAddress)

	monitoring.ServerStatusCreateChannel()
	go clientProcessGoroutine()

	for {
		tcpconn, err := mClientListener.AcceptTCP()
		if err != nil {
			utils.Logger.Warning("accept failed:", err)
			break
		}

		copyBackEndList := arraylist.New()
		copyBackendList(backEndList, copyBackEndList)

		//TODO: 최대 접속수 제한이 없음

		//TODO: Windows 에서는 버퍼 설정 안하는 것이 좋다
		tcpconn.SetReadBuffer(config.ClientSockReadbuf)
		tcpconn.SetWriteBuffer(config.ClientSockWritebuf)

		monitoring.ServerStatusCreateChannel()
		go handleClientTCPReceiveGoroutine(tcpconn, config, copyBackEndList)
	}

	utils.Logger.Info("tcpServerStart End")
}

func TCPServerEnd() {
	mClientListener.Close()
}

func checkError(err error) {
	if err != nil {
		utils.Logger.Fatal(err)
		os.Exit(-1)
	}
}

func getbackendServerList() *arraylist.List {
	backEndList := arraylist.New()
	backend := def.BackEndSessionForClient{}
	backEndList.Add(backend)
	// type BackEndSessionForClient struct {
	// 	UniqueIDNum int64
	// 	IP          net.IP
	// 	Conn        net.Conn
	// }

	return backEndList
}

func copyBackendList(sourceList *arraylist.List, targetList *arraylist.List) {
	for i := 0; i < sourceList.Size(); i++ {
		targetList.Add(sourceList.Get(i))
	}
	// sourceList.Each(func(index int, value interface{}) {
	// 	value, _ := sourceList.Get(index)
	// 	targetList.Add(value)
	// })
}
