package backEnd

import (
	"net"
	"sync/atomic"

	"../innerMsg"
	"../monitoring"
	packetDef "../packet"
	def "../typeDefine"
	"../utils"
)

var sequenceNum int32

func newSequenceIDNum() int32 {
	newValue := atomic.AddInt32(&sequenceNum, 1)
	return newValue
}

func connectGoroutine(backAddress string, config *def.Config) {
	utils.Logger.Info("start backEnd :", backAddress)

	defer monitoring.ServerStatusDestoryChannel()
	defer utils.PrintPanicStack()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", backAddress)
	if err != nil {
		utils.Logger.Error("fail ResolveTCPAddr backEnd :", err)
		return
	}

	var session def.BackEndSession
	session.InitSession()

	requestInnerMsgRegistSession(&session)

	for {
		if session.IsTryConnect() == false {
			break
		}

		backEnd, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			utils.Logger.Error("fail connect backEnd :", err)

			//TODO: 몇 초만 기다리기...

			continue
		}

		if ret := connectedSession(backEnd, &session); ret == false {
			continue
		}

		requestInnerMsgNewSession(&session)

		//TODO: Windows 에서는 버퍼 설정 안하는 것이 좋다
		backEnd.SetReadBuffer(config.ServerSockReadbuf)
		backEnd.SetWriteBuffer(config.ServerSockWritebuf)

		receiveBackEnd(&session, config)

		if session.IsTryConnect() {
			requestInnerMsgClosedSession(&session)
		}
	}

	utils.Logger.Error("ConnectGoroutine: wait end")

	session.WaitforEnd()

	utils.Logger.Error("ConnectGoroutine: end")
}

func receiveBackEnd(session *def.BackEndSession, config *def.Config) {
	utils.Logger.Info("start receiveBackEnd")

	maxPakcetSizeRd := config.ServerMaxPakcetSize
	HeaderSizeROnly := packetDef.HeaderSizeROnly
	var packetHeader packetDef.Header
	var readAbleByte int16
	var startRecvPos int16

	recviveBuff := make([]byte, maxPakcetSizeRd*3) //TODO: 매직넘버 수정해야 한다

	for {
		recvBytes, err := session.Conn.Read(recviveBuff[startRecvPos:])
		if recvBytes == 0 {
			//TODO: //requestInnerMsgConnectClose(session.UniqueIDNum)
			break
		}

		if err != nil {
			utils.Logger.Error("Tcp Read error: %s", err)
			//TODO: //requestInnerMsgConnectClose(session.UniqueIDNum)
			break
		}

		readAbleByte = startRecvPos + (int16)(recvBytes)
		var readPos int16

		for {
			if readAbleByte < HeaderSizeROnly {
				break
			}

			packetHeader.DecodingPacketHeader(recviveBuff[readPos:])
			requireDataSize := packetHeader.TotalSize

			if requireDataSize > readAbleByte {
				break
			}

			//TODO: 한번에 보내기로 한 패킷 보다 많이 보낸 경우
			if requireDataSize > maxPakcetSizeRd {
				utils.Logger.Warning("Larger than maximum send data: ", requireDataSize)
				break
			}

			ltvPacket := recviveBuff[readPos:(readPos + requireDataSize)]
			readPos += requireDataSize
			readAbleByte -= requireDataSize

			// 패킷 처리
			packetProcess(session, packetHeader.Id, requireDataSize, ltvPacket)
		}

		if readAbleByte > 0 {
			startRecvPos = readAbleByte
		}
	}

	utils.Logger.Debug("end receiveBackEnd")
}

func connectedSession(conn net.Conn, session *def.BackEndSession) bool {
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		utils.Logger.Error("[BackEndSession] cannot get remote address:", err)
		return false
	}

	session.UniqueIDNum = newSequenceIDNum()
	session.Conn = conn
	session.IP = net.ParseIP(host)

	utils.Logger.Info("[BackEndSession] new connection from:%v port:%v", host, port)
	return true
}

func requestInnerMsgRegistSession(session *def.BackEndSession) {
	msg := innerMsg.BackendServerData{}
	msg.MsgType = innerMsg.IdRegistBackEndSession
	msg.UniqueIDNum = session.UniqueIDNum
	msg.Session = session

	msgChannel := innerMsg.BackEndServerChannel()
	msgChannel <- msg
}

func requestInnerMsgNewSession(session *def.BackEndSession) {
	msg := innerMsg.BackendServerData{}
	msg.MsgType = innerMsg.IdNewBackEndSession
	msg.UniqueIDNum = session.UniqueIDNum

	msgChannel := innerMsg.BackEndServerChannel()
	msgChannel <- msg
}

func requestInnerMsgClosedSession(session *def.BackEndSession) {
	msg := innerMsg.BackendServerData{}
	msg.MsgType = innerMsg.IdClosedBackEndSession
	msg.UniqueIDNum = session.UniqueIDNum

	msgChannel := innerMsg.BackEndServerChannel()
	msgChannel <- msg
}

func packetProcess(session *def.BackEndSession, packetId int16, ltvSize int16, ltvPacket []byte) {
	if packetProcessIsRelay(session, packetId, ltvSize, ltvPacket) {
		return
	}

	//TODO:패킷 처리
}

func packetProcessIsRelay(session *def.BackEndSession, packetId int16, ltvSize int16, ltvPacket []byte) bool {
	if packetDef.IsRelayPacket(packetId) == false {
		return false
	}

	sendMsg := innerMsg.Backend2ClientMsg{}
	sendMsg.IsRelay = true
	sendMsg.LtvPacket = make([]byte, ltvSize)
	copy(ltvPacket, sendMsg.LtvPacket)

	inner := innerMsg.BackendProc2clientProcChannel()
	inner <- sendMsg
	return true
}
