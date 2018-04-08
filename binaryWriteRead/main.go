package main

import (
	"fmt"
)

// https://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
func EncodingLogIn(packetId int16, login *Login) ([]byte, int16) {
	totalSize := len(login.Id) + len(login.Pw) + HeaderSizeROnly
	packet := make([]byte, totalSize)

	writer := newWriter(packet, true)
	writer.WriteS16(totalSize)
	writer.WriteS16(packetId)
	writer.WriteBytes(login.Id)
	writer.WriteBytes(login.Pw)

	return packet, totalSize
}

func DecodingLogIn() {

}

// func (header *Header) DecodingPacketHeader(data []byte) {
// 	reader := newReader(data, true)
// 	header.TotalSize, _ = reader.ReadS16()
// 	header.Id, _ = reader.ReadS16()
// }

// func ProtocolEncodingPacket(packetid int16, bodySize int16, sendData []byte) ([]byte, int16) {
// 	totalSize := bodySize + HeaderSizeROnly
// 	sendBuf := make([]byte, totalSize)

// 	writer := newWriter(sendBuf, true)
// 	writer.WriteS16(totalSize)
// 	writer.WriteS16(packetid)
// 	writer.WriteBytes(sendData)
// 	return sendBuf, totalSize
// }

// func ProtocolEncodingToBackendRelayPacket(sessionId int64, packetId int16, ltvSize int16, ltvData []byte) []byte {
// 	totalSize := relayPacketHeaderSize + ltvSize
// 	sendBuf := make([]byte, totalSize)

// 	writer := newWriter(sendBuf, true)
// 	writer.WriteS16(totalSize)
// 	writer.WriteS16(packetId)
// 	writer.WriteS16((int16(1)))
// 	writer.WriteS64(sessionId)
// 	writer.WriteBytes(ltvData)
// 	return sendBuf
// }

func main() {
	fmt.Println("binary Write Read")

}
