package main

import (
	"net"
	"os"
)

const (
	RECV_BUF_LEN = 1024
)

func main() {
	if len(os.Args) == 1 {
		println("need request parameter")
		os.Exit(1)
	}
	echo_contents := os.Args[1]
	tcp_addr, err := net.ResolveTCPAddr("tcp", "localhost:6666")
	if err != nil {
		println("error tcp resolve failed", err.Error())
		os.Exit(1)
	}
	tcp_conn, err := net.DialTCP("tcp", nil, tcp_addr)
	SendEcho(tcp_conn, echo_contents)

	echo := GetEcho(tcp_conn)
	println("echo: ", string(echo))
	println("receive success")
	tcp_conn.Close()
}

func SendEcho(conn *net.TCPConn, msg string) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		println("Error send request:", err.Error())
	} else {
		println("Request sent")
	}
}

func GetEcho(conn *net.TCPConn) string {
	buf_recever := make([]byte, RECV_BUF_LEN)
	_, err := conn.Read(buf_recever)
	if err != nil {
		println("Error while receive response:", err.Error())
		return ""
	}
	return string(buf_recever)
}
  
