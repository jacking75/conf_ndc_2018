package packet

import (
	"fmt"
	"reflect"
)

const (
	ProtocolidFirst = 1

	ProtocolidDevFirst   = 11
	ProtocolidDevEchoReq = 14
	ProtocolidDevEchoRes = 15
	ProtocolidDevEnd     = 21

	ProtocolidEnd = 500
)

type Header struct {
	TotalSize int16
	Id        int16
}

var HeaderSizeROnly = protocolInitHeaderSize()

const relayPacketHeaderSize int16 = 12 // TotalSize(2)+패킷아이디(2)+대상수(2)+[대상(8)]

func (header *Header) DecodingPacketHeader(data []byte) {
	reader := newReader(data, true)
	header.TotalSize, _ = reader.ReadS16()
	header.Id, _ = reader.ReadS16()
}

func ProtocolEncodingPacket(packetid int16, bodySize int16, sendData []byte) ([]byte, int16) {
	totalSize := bodySize + HeaderSizeROnly
	sendBuf := make([]byte, totalSize)

	writer := newWriter(sendBuf, true)
	writer.WriteS16(totalSize)
	writer.WriteS16(packetid)
	writer.WriteBytes(sendData)
	return sendBuf, totalSize
}

func ProtocolEncodingToBackendRelayPacket(sessionId int64, packetId int16, ltvSize int16, ltvData []byte) []byte {
	totalSize := relayPacketHeaderSize + ltvSize
	sendBuf := make([]byte, totalSize)

	writer := newWriter(sendBuf, true)
	writer.WriteS16(totalSize)
	writer.WriteS16(packetId)
	writer.WriteS16((int16(1)))
	writer.WriteS64(sessionId)
	writer.WriteBytes(ltvData)
	return sendBuf
}

type SubsetData struct {
	Id   int16
	Data []byte
}

func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		fmt.Println("reflect.Array")
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		fmt.Println("reflect.Struct")
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		//fmt.Println("reflect.int")
		return int(t.Size())
	case reflect.Slice:
		//fmt.Println("reflect.Slice:", sizeof(t.Elem()))
		return 0
	}

	return -1

}

func protocolInitHeaderSize() int16 {
	var packetHeader Header
	headerSize := sizeof(reflect.TypeOf(packetHeader))
	return (int16)(headerSize)
}

func IsRelayPacket(packetId int16) bool {
	if packetId > ProtocolidEnd {
		return true
	}
	return false
}
