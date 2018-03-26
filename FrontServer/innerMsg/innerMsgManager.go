package innerMsg

type innerMsgManager struct {
	clientMsgChannel        chan ClientData
	backendServerMsgChannel chan BackendServerData

	backendProc2clientProcChannel chan Backend2ClientMsg
	clientProc2backendProcChannel chan Client2BackendMsg
}

var innerMsgManagerInstance = &innerMsgManager{}

func InnerMsgManagerInit(maxBackendServerChannelSize int, maxClientChannelSize int) {
	innerMsgManagerInstance.clientMsgChannel = make(chan ClientData, maxClientChannelSize)
	innerMsgManagerInstance.backendServerMsgChannel = make(chan BackendServerData, maxBackendServerChannelSize)

	//TODO: 채널 크기 인자로 받기
	innerMsgManagerInstance.backendProc2clientProcChannel = make(chan Backend2ClientMsg, 100)
	innerMsgManagerInstance.clientProc2backendProcChannel = make(chan Client2BackendMsg, 100)
}

func ClientChannel() chan ClientData {
	return innerMsgManagerInstance.clientMsgChannel
}

func BackEndServerChannel() chan BackendServerData {
	return innerMsgManagerInstance.backendServerMsgChannel
}

func BackendProc2clientProcChannel() chan Backend2ClientMsg {
	return innerMsgManagerInstance.backendProc2clientProcChannel
}

func ClientProc2backendProcChannel() chan Client2BackendMsg {
	return innerMsgManagerInstance.clientProc2backendProcChannel
}

/*
type Client2BackendMsg struct {
	isRay             bool
	ltvPacket         []byte
}

type Backend2ClientMsg struct {
	UniqueIDNum int64
	ltvPacket   []byte
}
*/
