package main

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"sync"
)

var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)

func main() {
	// 监听在本地7890端口
	listenAddr := "127.0.0.1:7890"
	proxyAddr := "127.0.0.1:12345" // 这里是你的代理服务器地址和端口

	// 创建一个TCP socket，用于监听浏览器的HTTP请求
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v\n", listenAddr, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Proxy server listening on %s\n", listenAddr)

	for {
		// 等待浏览器的连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept browser connection: %v\n", err)
			continue
		}

		// 在新的goroutine中处理浏览器请求
		red.Println("A new browser request:")
		go handleBrowserRequest(conn, proxyAddr)
	}
}

func handleBrowserRequest(browserConn net.Conn, proxyAddr string) {
	var wg sync.WaitGroup

	defer func(browserConn net.Conn) {
		err := browserConn.Close()
		if err != nil {

		}
	}(browserConn)

	// 连接到代理服务器
	proxyServerConn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		fmt.Printf("Failed to connect to proxy server: %v\n", err)
		return
	}
	defer func(proxyServerConn net.Conn) {
		err := proxyServerConn.Close()
		if err != nil {

		}
	}(proxyServerConn)

	// 转发数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 1024)
		for {
			n, err := browserConn.Read(buf)
			if err != nil {
				fmt.Printf("Error reading from browserConn: %v\n", err)
				break
			}
			n, err = green.Print(string(buf[:n]))
			if err != nil {
				fmt.Println(err, n)
			}
			_, err = proxyServerConn.Write(buf[:n])
			if err != nil {
				fmt.Printf("Failed to forward data to proxy server: %v\n", err)
				break
			}

			response := []byte("HTTP/1.1 200 OK\r\n" +
				"Content-Type: text/html\r\n" +
				"Content-Length: 69\r\n" +
				"\r\n" +
				"<html><body><h1>Hello, World! My name is LosFurina</h1></body></html>")

			_, err = browserConn.Write(response)
			if err != nil {
				fmt.Printf("Error writing response to browserConn: %v\n", err)
			}

		}
	}()
	wg.Wait()

	//go func() {
	//	buf := make([]byte, 1024)
	//	for {
	//		n, err := proxyServerConn.Read(buf)
	//		if err != nil {
	//			fmt.Printf("Error reading from proxyServerConn: %v\n", err)
	//			break
	//		}
	//		n, err = red.Print(string(buf[:n]))
	//		if err != nil {
	//			fmt.Println(err, n)
	//		}
	//		_, err = browserConn.Write(buf[:n])
	//		if err != nil {
	//			fmt.Printf("Failed to forward data to browser: %v\n", err)
	//			break
	//		}
	//	}
	//}()
}
