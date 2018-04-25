package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	RECV_BUF_LEN = 1024
)

func main() {

	var msg string
	var port int

	flag.IntVar(&port, "p", 32452, "port")
	flag.StringVar(&msg, "m", "", "send message")
	flag.Parse()

	if len(msg) == 0 {
		println("send message is empty")
		os.Exit(1)
	}

	tcp_addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		println("error tcp resolve failed", err.Error())
		os.Exit(1)
	}

	tcp_conn, err := net.DialTCP("tcp", nil, tcp_addr)
	if err != nil {
		println("error tcp conn failed", err.Error())
		os.Exit(1)
	}
	SendEcho(tcp_conn, msg)

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
