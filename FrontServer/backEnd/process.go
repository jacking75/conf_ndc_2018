package backEnd

import (
	"../innerMsg"
	"../monitoring"
	def "../typeDefine"
	"../utils"
)

func processGoroutine() {
	utils.Logger.Info("start BackEnd processGoroutine")

	defer monitoring.ServerStatusDestoryChannel()
	defer utils.PrintPanicStack()

	serverDic := make(map[int32]*def.BackEndSession)
	msgChannel := innerMsg.BackEndServerChannel()
	terminateNotifyChannel := utils.GetServerTerminateChennel()
	clientProc2backendProcChannel := innerMsg.ClientProc2backendProcChannel()

loop:
	for {
		select {
		case msg := <-msgChannel:
			{
				switch msg.MsgType {
				case innerMsg.IdRegistBackEndSession:
					{
						//TODO 실패에 대한 처리 필요
						addSession(msg.Session, serverDic)
					}
				case innerMsg.IdClosedBackEndSession:
					{
						// ?
						DecrementConnectCount()
					}
				case innerMsg.IdNewBackEndSession:
					{
						IncrementConnectCount()
					}
				case innerMsg.IdRemoveBackEndSession:
					{
						//?
					}
				}
			}
		case client2BackendMsg := <-clientProc2backendProcChannel:
			{
				if client2BackendMsg.IsRelay {
					//TODO: 프론트가 백엔드에 요청을 보낼 때. 예) 클라이언트 인증, 클라이언트 고유 번호 통보
				} else {
					relayToBackend(client2BackendMsg)
				}
			}
		case <-terminateNotifyChannel:
			break loop
		}
	}

	utils.Logger.Debug("pre end processGoroutine")

	allRemoveSession(serverDic)

	utils.Logger.Debug("end processGoroutine")
}

func addSession(session *def.BackEndSession, dic map[int32]*def.BackEndSession) bool {
	dic[session.UniqueIDNum] = session
	return true
}

func removeSession(sessionIDNum int32, dic map[int32]*def.BackEndSession) bool {
	if session, ok := dic[sessionIDNum]; ok {
		session.DisableTryConnect()
		session.SetEnd()
		session.TCPClose()

		delete(dic, sessionIDNum)
		return true
	}

	return false
}

func allRemoveSession(dic map[int32]*def.BackEndSession) {
	utils.Logger.Info("start allRemoveSession")

	for key := range dic {
		removeSession(key, dic)
	}

	utils.Logger.Debug("end allRemoveSession")
}

func relayToBackend(clientMsg innerMsg.Client2BackendMsg) {
	// TODO: 여기에 이 메시지를 보낸 유저가 어떤 백엔드에 연결되어 있는지 아는 정보가 있어야 한다.
}
