package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"time"
	"io"
	"net"
	"os"

	quic "github.com/lucas-clemente/quic-go"
)

func receiveData(conn *net.UDPConn) {
	for {
		fmt.Println("waiting for data here\n")
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
	}
}

func handleData(stream quic.Stream) {
	buf := make([]byte, 1024)
	for {
		fmt.Println("Client: waiting for data\n")
		_, err := io.ReadFull(stream, buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Client: Got '%s'\n", buf)
	}
}

const host = "172.20.114.49"
const port = 9981
const message = "foobar"
const addr = "172.20.114.49:9982"

func main() {
	keyLog, err := os.Create("clientkey.log")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("create key log: %s\n", keyLog)
	if err != nil {
		fmt.Printf("Could not create key log: %s\n", err.Error())
		os.Exit(1)
	}
	if keyLog != nil {
		defer keyLog.Close()
	}
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
		KeyLogWriter:       keyLog,
	}

	quicConf := &quic.Config{
		HandshakeTimeout: 24*time.Hour,
		MaxIdleTimeout:   24*time.Hour,
	}
	session, err := quic.DialAddr(addr, tlsConf, quicConf)
	if err != nil {
		fmt.Println(err)
		return
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("Client: Sending '%s'\n", message)
	// _, err = stream.Write([]byte(message))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	go handleData(stream)

	// ip := net.ParseIP(host)
	// srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	// dstAddr := &net.UDPAddr{IP: ip, Port: port}
	// conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer conn.Close()
	// // conn.Write([]byte("hello"))
	// fmt.Printf("<%s>\n", conn.RemoteAddr())

	// // go func(conn *net.UDPConn) {
	// // 	for {
	// // 		fmt.Printf("waiting for data here\n")
	// // 		data := make([]byte, 1024)
	// // 		n, err := conn.Read(data)
	// // 		fmt.Printf("read %s from <%s>\n", data[:n], conn.RemoteAddr())
	// // 	}
	// // }()

	// go receiveData(conn)

	var (
		input string
		cmd   string
		param string
	)

	f := bufio.NewReader(os.Stdin) //读取输入的内容
	for {
		fmt.Println("请输入一些字符串>")
		input, _ = f.ReadString('\n') //定义一行输入的内容分隔符。
		if len(input) == 1 {
			continue //如果用户输入的是一个空行就让用户继续输入。
		}
		fmt.Printf("您输入的是:%s", input)
		fmt.Sscan(input, &cmd, &param)
		switch cmd {
		case "w":
			// conn.Write([]byte(param))
			fmt.Printf("===============Client: Sending '%s'\n", param)
			var sendStr string
			for i := 0; i < 1024; i++ {
				sendStr += param
			}
			_, err = stream.Write([]byte(sendStr))
			if err != nil {
				fmt.Println(err)
				return
			}
		case "stop":
			break
		default:
			fmt.Printf("未知命令：%s\n", cmd)
		}
		fmt.Printf("您输入的第一个参数是 %s, 输入的第二个参数是 %s\n", cmd, param)
	}
}
