package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
)

var green = color.New(color.FgGreen)

func main() {
	listenAddress := "127.0.0.1:12345" // 服务器监听地址，可根据需要修改

	// 创建一个TCP监听器
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v\n", listenAddress, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s\n", listenAddress)

	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept client connection: %v\n", err)
			continue
		}

		// 启动一个新goroutine来处理客户端连接
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Close connection failed!")
		}
	}(conn)

	// 读取客户端发送的数据
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("Client connection closed: %v\n", err)
			break
		}

		data := buffer[:n]
		n, err = green.Print(string(data))
		if err != nil {
			fmt.Println(err, n)
		}
		// 在这里可以对接收到的数据进行处理
		// 例如，你可以将数据存储到文件、进行解析、发送响应等操作
	}
}
