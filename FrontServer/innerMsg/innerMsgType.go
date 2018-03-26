package innerMsg

import (
	def "../typeDefine"
)

const (
	// client side
	IdClientNewSession       = 101
	IdClientRemoveSession    = 102
	IdClientRemoveAllSession = 103

	// backend side
	IdRegistBackEndSession = 201
	IdNewBackEndSession    = 202
	IdClosedBackEndSession = 203

	IdRemoveBackEndSession    = 286
	IdRemoveAllBackEndSession = 287

	// client-backend side
	IdRelayDataB2C = 301
	//IdNtfNewBackendInfoB2C = 301
)

type ClientData struct {
	UniqueIDNum int64
	Session     *def.Session
	MsgType     int16
	Data        []byte
}

type BackendServerData struct {
	UniqueIDNum int32
	Session     *def.BackEndSession
	MsgType     int16
	Data        []byte
}

type Client2BackendMsg struct {
	ClientUniqueIDNum int64
	IsRelay           bool
	LtvPacket         []byte
}

type Backend2ClientMsg struct {
	IsRelay   bool
	LtvPacket []byte
}
