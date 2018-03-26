package clientGateway

import (
	"net"
	_ "os"
	"sync/atomic"

	"../innerMsg"
	"../monitoring"
	packetDef "../packet"
	def "../typeDefine"
	"../utils"

	"github.com/emirpasic/gods/lists/arraylist"
)

var sequenceNum int64

func newSequenceIDNum() int64 {
	newValue := atomic.AddInt64(&sequenceNum, 1)
	return newValue
}

func handleClientTCPReceiveGoroutine(conn net.Conn, config *def.Config, backEndList *arraylist.List) {
	logger := utils.Logger

	defer monitoring.ServerStatusDestoryChannel()
	defer utils.PrintPanicStack()

	//TODO: Session 생성은 syncpool 사용으로 변경하자
	var session def.Session
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		logger.Error("cannot get remote address:", err)
		return
	}

	session.UniqueIDNum = newSequenceIDNum()
	session.Conn = conn
	session.IP = net.ParseIP(host)

	logger.Info("new connection from:%v port:%v", host, port)

	//TODO: 버퍼 크기를 크게 한다.
	in := make(chan []byte)
	defer func() {
		close(in) // session will close
	}()

	// session die signal, will be triggered by agent()
	//session.Die = make(chan struct{}, 3)
	session.WriteEnd = make(chan struct{}, 1)
	session.ReceiveEnd = make(chan struct{}, 1)
	session.TCPSendChannel = make(chan def.TCPSendData, config.MaxClientTCPSendChannelSize)

	requestInnerMsgNewSession(&session)

	// 이 세션의 send용 고루틴 생성한다
	monitoring.ServerStatusCreateChannel()
	go handleClientTCPSendGoroutine(&session, config)

	maxPakcetSizeRd := config.ClientMaxPakcetSize
	HeaderSizeROnly := packetDef.HeaderSizeROnly
	var packetHeader packetDef.Header
	var readAbleByte int16
	var startRecvPos int16
	recviveBuff := make([]byte, maxPakcetSizeRd*3) //TODO: 매직넘버 수정해야 한다

	for {
		recvBytes, err := conn.Read(recviveBuff[startRecvPos:])
		if recvBytes == 0 {
			requestInnerMsgConnectClose(session.UniqueIDNum)
			break
		}

		if err != nil {
			logger.Error("Tcp Read error: %s", err)

			requestInnerMsgConnectClose(session.UniqueIDNum)
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
				logger.Warning("Larger than maximum send data: ", requireDataSize)
				break
			}

			ltvPacket := recviveBuff[readPos:(readPos + requireDataSize)]
			readPos += requireDataSize
			readAbleByte -= requireDataSize

			// 패킷 처리
			packetProcess(&session, packetHeader.Id, requireDataSize, ltvPacket)
		}

		if readAbleByte > 0 {
			startRecvPos = readAbleByte
		}
	}

	utils.Logger.Debug("wait session.ReceiveEnd")
	select {
	case <-session.ReceiveEnd:
		//case  //TODO: 지정 시간까지 통보가 안오면 자동으로 고루틴에서 나가도록 한다
	}

	utils.Logger.Debug("end handleClientTCPReceiveGoroutine")
}

func requestInnerMsgNewSession(session *def.Session) {
	msg := innerMsg.ClientData{}
	msg.MsgType = innerMsg.IdClientNewSession
	msg.UniqueIDNum = session.UniqueIDNum
	msg.Session = session

	msgClientChannel := innerMsg.ClientChannel()
	msgClientChannel <- msg
}

func requestInnerMsgConnectClose(sessionUnique int64) {
	msg := innerMsg.ClientData{UniqueIDNum: sessionUnique, MsgType: innerMsg.IdClientRemoveSession}

	msgChannel := innerMsg.ClientChannel()
	msgChannel <- msg
}

func packetProcess(session *def.Session, packetId int16, ltvSize int16, ltvPacket []byte) {
	if packetProcessIsRelay(session, packetId, ltvSize, ltvPacket) {
		return
	}

	if isProcess := packetProcessImmediately(session, packetId, ltvSize, ltvPacket); isProcess == false {
		//TODO: 다른 서버로 릴레이 하거나 clientProcess에서 처리하도록 한다
	}
}

func packetProcessIsRelay(session *def.Session, packetId int16, ltvSize int16, ltvPacket []byte) bool {
	if packetDef.IsRelayPacket(packetId) == false {
		return false
	}

	sendMsg := innerMsg.Client2BackendMsg{}
	sendMsg.IsRelay = true
	var sessionId int64 = 1 //TODO: 세션아이디 지정하기
	sendMsg.LtvPacket = packetDef.ProtocolEncodingToBackendRelayPacket(sessionId, packetId, ltvSize, ltvPacket)

	inner := innerMsg.ClientProc2backendProcChannel()
	inner <- sendMsg
	return true
}

func packetProcessImmediately(session *def.Session, packetId int16, ltvSize int16, ltvPacket []byte) bool {
	//TODO: 클라이언트가 접속 중에 백엔드 서버 정보가 바뀌는 경우 여기에서 정보를 받을 수 있도록 하자
	// 여기에서만 사용하는 채널이 있으면 좋을 듯. atomic 변수(처리할 메시지 숫자. 입력, 처리마다 숫자를 변동)로
	// 채널을 읽도록 하면 읽고 처리한 후 완료한다
	isCompleted := packetProcessDev(session, packetId, ltvSize, ltvPacket)

	return isCompleted
}

func packetProcessDev(session *def.Session, packetId int16, ltvSize int16, ltvPacket []byte) bool {
	if packetId < packetDef.ProtocolidDevFirst || packetId > packetDef.ProtocolidDevEnd {
		return false
	}

	isSocketNoProblem := true
	bodySize := ltvSize - packetDef.HeaderSizeROnly
	bodyData := ltvPacket[packetDef.HeaderSizeROnly:]

	switch packetId {
	case packetDef.ProtocolidDevEchoReq:
		sendData, dataSize := packetDef.ProtocolEncodingPacket(packetDef.ProtocolidDevEchoRes, bodySize, bodyData)
		isSocketNoProblem = packetSendImmediately(session.Conn, dataSize, sendData)
	}

	if isSocketNoProblem == false {
		//TODO: 소켓을 짜르도록 한다
		//TODO: 못 보낸 데이터를 다음에 보내고 싶다면 어딘가에 저장해야 한다
	}

	return true
}

///// 데이터 보내기
func packetSendImmediately(conn net.Conn, sendDataSize int16, sendData []byte) bool {
	completeSize, ret := conn.Write(sendData)

	if ret != nil {
		utils.Logger.Warning("Error packetSendImmediately, bytes: %v reason: %v", completeSize, ret)
		return false
	}

	return true
}
