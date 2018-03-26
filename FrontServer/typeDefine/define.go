package typeDefine

import (
	"net"
	"time"
)

const (
	SessionFlagKeyExcg    = 0x1 // KEY 교환 완료
	SessionFlagEnCrypt    = 0x2 // 암호화 사용
	SessionFlagForceClose = 0x4 // 쫒아내기
	SessionFlagAuthorized = 0x8 // 인증 완료
)

type Config struct {
	MachineId int64 // snowflake를 위한 유니크 id

	// Client Side
	ClientAddress string // 만약 IP와 포트번호 결합이면 localhost:19999

	//readDeadline time.Duration
	ClientSockReadbuf   int
	ClientSockWritebuf  int
	ClientMaxPakcetSize int16

	ServerSockReadbuf   int
	ServerSockWritebuf  int
	ServerMaxPakcetSize int16

	MaxServerInnerMsgChannelSize, MaxClientInnerMsgChannelSize int
	MaxServerTCPSendChannelSize, MaxClientTCPSendChannelSize   int

	Nodelay int

	Txqueuelen int // ?
	Dscp       int // ?
	Interval   int // ?
	Resend     int // ?
	NC         int // ?

	// Server Side
}

func (config *Config) Setting() {
	config.MachineId = 77
	config.ClientAddress = "localhost:32452"
	//config.readDeadline: ,
	config.ClientSockReadbuf = 8000
	config.ClientSockWritebuf = 8000
	config.ClientMaxPakcetSize = 1024
	config.Nodelay = 1

	config.ServerSockReadbuf = 16000
	config.ServerSockWritebuf = 16000
	config.ServerMaxPakcetSize = 4028

	config.MaxServerInnerMsgChannelSize = 1000
	config.MaxClientInnerMsgChannelSize = 1000

	config.MaxServerTCPSendChannelSize = 128
	config.MaxClientTCPSendChannelSize = 64

	//txqueuelen:   c.Int("txqueuelen"),
	//dscp:         c.Int("dscp"),
	//interval:     c.Int("interval"),
	//resend:       c.Int("resend"),
	//nc:           c.Int("nc"),
}

type Session struct {
	UniqueIDNum int64
	IP          net.IP
	Conn        net.Conn

	//MQ      chan pb.Game_Frame          // 返回给客户端的异步消息
	//Encoder *rc4.Cipher                 // 加密器
	//Decoder *rc4.Cipher                 // 解密器
	//UserId int32 // 玩家ID
	//GSID    string                      // 游戏服ID;e.g.: game1,game2
	//Stream  pb.GameService_StreamClient // 后端游戏服数据流
	//Die chan struct{}

	WriteEnd   chan struct{}
	ReceiveEnd chan struct{}

	TCPSendChannel chan TCPSendData

	//TODO: 프로세스가 보내는 메시지를 받는 채널
	//ProcessToChannel chan interface{}

	// 会话标记
	Flag int32

	// 时间相关
	ConnectTime    time.Time // TCP链接建立时间
	PacketTime     time.Time // 当前包的到达时间
	LastPacketTime time.Time // 前一个包到达时间

	PacketCount     uint32 // 对收到的包进行计数，避免恶意发包
	PacketCount1Min int    // 每分钟的包统计，用于RPM判断
}

type TCPSendData struct {
	SessionUniqueIDNum int64
	Data               []byte
}

type BackEndSession struct {
	UniqueIDNum int32
	IP          net.IP
	Conn        net.Conn

	isTryConnect bool
	waitForEnd   chan struct{}
}

func (session *BackEndSession) InitSession() {
	session.isTryConnect = true
	session.waitForEnd = make(chan struct{}, 1)
}

func (session *BackEndSession) IsTryConnect() bool {
	if session.isTryConnect == true {
		return true
	}

	return false
}

func (session *BackEndSession) TCPClose() {
	session.Conn.Close()
}

func (session *BackEndSession) DisableTryConnect() {
	session.isTryConnect = false
}

func (session *BackEndSession) SetEnd() {
	close(session.waitForEnd)
}

func (session *BackEndSession) WaitforEnd() {
	select {
	case <-session.waitForEnd:
		//case  //TODO: 지정 시간까지 통보가 안오면 자동으로 고루틴에서 나가도록 한다
	}
}

type BackEndSessionForClient struct {
	UniqueIDNum int64
	Conn        net.Conn
}

type BackEndServerConnectInfo struct {
	Address string // 만약 IP와 포트번호 결합이면 localhost:19999
	// Server Side
}
