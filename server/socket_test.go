package server

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
)

func TestSocketServer(t *testing.T) {
	fmt.Println("start server .....")

	listen, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		fmt.Println("listen failed...", err)
		return
	}

	for {
		conn, err := listen.Accept()

		if err != nil {
			fmt.Println("listen failed...", err)
			return
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 100)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("read data size %d msg:%s", n, string(buf[0:n]))
		msg := []byte("hello,world\n")
		conn.Write(msg)
	}
}

func TestSocketClient(t *testing.T) {
	fmt.Println("start client ......")

	conn, err := net.Dial("tcp", "0.0.0.0:8080")

	if err != nil {
		fmt.Println("dial failed...", err)
		return
	}

	defer conn.Close()

	inputReader := bufio.NewReader(os.Stdin)

	for {
		str, _ := inputReader.ReadString('\n')
		data := strings.Trim(str, "\n")

		//data := "jinyidong"

		if data == "quit" { //输入quit退出
			return
		}
		_, err := conn.Write([]byte(data)) //发送数据
		if err != nil {
			fmt.Println("send data error:", err)
			return
		}
		buf := make([]byte, 512)
		n, err := conn.Read(buf) //读取服务端端数据
		fmt.Println("from server:", string(buf[:n]))
	}
}
