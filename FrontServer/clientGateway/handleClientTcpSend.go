package clientGateway

import (
	"../monitoring"
	def "../typeDefine"
	"../utils"
)

func handleClientTCPSendGoroutine(session *def.Session, config *def.Config) {
	utils.Logger.Debug("start handleClientTcpSendGoroutine")

	defer monitoring.ServerStatusDestoryChannel()
	defer utils.PrintPanicStack()

	writeEnd := session.WriteEnd
	//tcpSendChannel := session.tcpSendChannel

loop:
	for {
		select {
		//case sendData := <-tcpSendChannel:
		//	{
		//TODO: 클라이언트에세 데이터를 보낸다
		//	}
		case <-writeEnd:
			break loop
		}
	}

	utils.Logger.Debug("end handleClientTcpSendGoroutine")
	close(session.ReceiveEnd)
}
