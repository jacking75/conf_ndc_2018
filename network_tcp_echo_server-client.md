## cho Server/Client
echoClient.go    
```Go
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
```  
  
echoServer.go    
```Go
package main

import (
  "fmt"
  "net"
	"os"
)

const (
	RECV_BUF_LEN = 1024
)

func main() {
	fmt.Println("Starting the server")

	listener, err := net.Listen("tcp", "localhost:6666")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accept:", err.Error())
			return
		}
		go EchoFunc(conn)
	}
}

func EchoFunc(conn net.Conn) {
	buf := make([]byte, RECV_BUF_LEN)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:%d, %s", n, err.Error())
		return
	}
	fmt.Println("received ", n, " bytes of data =", string(buf))

	//send reply
	_, err = conn.Write(buf)
	if err != nil {
		fmt.Println("Error send reply:", err.Error())
	} else {
		fmt.Println("Reply sent")
	}
        conn.Close()
}
```