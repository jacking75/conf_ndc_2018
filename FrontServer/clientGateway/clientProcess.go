package clientGateway

import (
	"encoding/binary"
	"net"

	"../innerMsg"
	"../monitoring"
	def "../typeDefine"
	"../utils"
)

func clientProcessGoroutine() {
	utils.Logger.Debug("start clientProcessGoroutine")

	defer monitoring.ServerStatusDestoryChannel()
	defer utils.PrintPanicStack()

	innerMsgLoop()

	utils.Logger.Debug("end clientProcessGoroutine")
}

func innerMsgLoop() {
	clientDic := make(map[int64]*def.Session)
	msgChannel := innerMsg.ClientChannel()
	backendProc2clientProcChannel := innerMsg.BackendProc2clientProcChannel()
	terminateNotifyChannel := utils.GetServerTerminateChennel()

	defer allRemoveSession(clientDic)

loop:
	for {
		select {
		case msg := <-msgChannel:
			{
				//TODO: 각 메시지 처리
				switch msg.MsgType {
				case innerMsg.IdClientNewSession:
					{
						//TODO 실패에 대한 처리 필요
						addSession(msg.Session, clientDic)
					}
				case innerMsg.IdClientRemoveSession:
					{
						removeSession(msg.UniqueIDNum, clientDic)
					}
				case innerMsg.IdClientRemoveAllSession:
				}
			}
		case backend2ClientMsg := <-backendProc2clientProcChannel:
			{
				processInnerMsgFromBackend(clientDic, &backend2ClientMsg)
			}
		case <-terminateNotifyChannel:
			break loop
		}

		//if session.Flag&sessionFlagForceClose != 0 {
		//	return
		//}
	}
}
func addSession(session *def.Session, dic map[int64]*def.Session) bool {
	dic[session.UniqueIDNum] = session
	return true
}

func removeSession(sessionIDNum int64, dic map[int64]*def.Session) bool {
	if session, ok := dic[sessionIDNum]; ok {
		tcpClose(session.Conn)
		close(session.WriteEnd)

		delete(dic, sessionIDNum)
		return true
	}

	return false
}

func allRemoveSession(dic map[int64]*def.Session) {
	utils.Logger.Info("start allRemoveSession")

	for key := range dic {
		removeSession(key, dic)
	}

	utils.Logger.Debug("end allRemoveSession")
}

func findSession(sessionIDNum int64, dic map[int64]*def.Session) *def.Session {
	if session, ok := dic[sessionIDNum]; ok {
		return session
	}
	return nil
}

func processInnerMsgFromBackend(dic map[int64]*def.Session, msg *innerMsg.Backend2ClientMsg) {
	if processInnerMsgRelayPacket(dic, msg) {
		return
	}

	//TODO: 백엔드에서 보낸 메시지를 처리한다.
}

func processInnerMsgRelayPacket(dic map[int64]*def.Session, msg *innerMsg.Backend2ClientMsg) bool {
	if msg.IsRelay == false {
		return false
	}

	var clientCount uint16
	binary.LittleEndian.PutUint16(msg.LtvPacket[3:], clientCount)

	ltvStartPos := 6 + (clientCount * 8)
	ltvPacket := msg.LtvPacket[ltvStartPos:]

	var clientPos int = 6
	for i := (uint16)(0); i < clientCount; i++ {
		var clientUniqueId uint64
		binary.LittleEndian.PutUint64(msg.LtvPacket[clientPos:], clientUniqueId)

		if session := findSession((int64)(clientUniqueId), dic); session != nil {
			sendData := def.TCPSendData{}
			sendData.SessionUniqueIDNum = (int64)(clientUniqueId)
			sendData.Data = ltvPacket

			session.TCPSendChannel <- sendData
		} else {
			//TODO: 로그 남겨야 하나?
		}
	}
	return true
}

type TCPSendData struct {
	SessionUniqueIDNum uint64
	Data               []byte
}

//func clientProcess(session *sessionSt, in chan []byte) {
//logger := utils.Logger

//defer gVarServerStatus.serverStatusDestoryChannel()
//defer gVarWait.Done() // will decrease waitgroup by one, useful for manual server shutdown
//defer utils.PrintPanicStack()

// init session
//sess.MQ = make(chan pb.Game_Frame, 512)

// minute timer
//min_timer := time.After(time.Minute)

// cleanup work
//defer func() {
//tcpClose(session.Conn)

//close(session.Die)
//if sess.Stream != nil {
//	sess.Stream.CloseSend()
//}
//}()

// >> the main message loop <<
// handles 4 types of message:
//  1. from client
//  2. from game service
//  3. timer
//  4. server shutdown signal
//for {
//	select {

/*case msg, ok := <-in: // packet from network
if !ok {
	return
}

session.ConnectTime = time.Now()
session.LastPacketTime = time.Now()
session.PacketCount++
session.PacketCount1Min++
session.PacketTime = time.Now()

if result := route(sess, msg); result != nil {
	out.send(sess, result)
}
sess.LastPacketTime = sess.PacketTime
*/

//case frame := <-sess.MQ: // packets from game
//	switch frame.Type {
//	case pb.Game_Message:
//		out.send(sess, frame.Message)
//	case pb.Game_Kick:
//		sess.Flag |= SESS_KICKED_OUT
//	}
//case <-min_timer: // minutes timer
//	timer_work(sess, out)
//	min_timer = time.After(time.Minute)

//case <-session.Die:
//	logger.Info("closed client session")
//	return
//case <-gVarServerTerminate: // server is shuting down...
//	session.Flag |= sessionFlagForceClose
//}

// see if the player should be kicked out.
//if session.Flag&sessionFlagForceClose != 0 {
//	return
//}
//}
//}

func tcpClose(tcp net.Conn) {
	tcp.Close()
	//TODO: close 2번 호출해도 문제 없는지 테스트 해본다
	//tcp.Close()
}
