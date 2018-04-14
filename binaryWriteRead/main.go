package main

import (
	"fmt"
)

func (login *Login) Encoding() ([]byte, int16) {
	totalSize := (int16)(len(login.Id)) + (int16)(len(login.Pw)) + HeaderSizeROnly
	packet := make([]byte, totalSize)

	writer := newWriter(packet, true)
	writer.WriteS16(totalSize)
	writer.WriteS16(login.PktId)
	writer.WriteBytes(login.Id[:], 16)
	writer.WriteBytes(login.Pw[:], 16)

	return packet, totalSize
}

func (login *Login) DecodingLogIn(data []byte, size int16) {
	reader := newReader(data, true)
	login.TotalSize, _ = reader.ReadS16()
	login.PktId, _ = reader.ReadS16()
	id, _ := reader.ReadBytes(16)
	pw, _ := reader.ReadBytes(16)

	copy(login.Id[:], id[:])
	copy(login.Pw[:], pw[:])
}

func main() {
	fmt.Println("binary Write Read")

	id := "jacking"
	pw := "123ert"

	var pktLoging Login
	pktLoging.PktId = int16(11)
	copy(pktLoging.Id[:], []byte(id))
	copy(pktLoging.Pw[:], []byte(pw))

	pktData, pktSize := (&pktLoging).Encoding()
	fmt.Println("Login packet Encoding Info")
	fmt.Println(pktSize)
	fmt.Println("PktId: ", pktLoging.PktId)
	fmt.Println("Id: ", pktLoging.Id)
	fmt.Println("Pw: ", pktLoging.Pw)

	fmt.Println("")

	var pktLoging2 Login
	(&pktLoging2).DecodingLogIn(pktData, pktSize)
	fmt.Println("Login packet decoding Info")
	fmt.Println("PktId: ", pktLoging2.PktId)
	fmt.Println("Id: ", pktLoging2.Id)
	fmt.Println("Pw: ", pktLoging2.Pw)

	fmt.Println("sting Id: ", string(pktLoging2.Id[:]))
}
