package main

import (
	"flag"
	"fmt"
	"io"
	"msghandler"
	"net"
	"protocol"
	"syscall"
)

func main() {
	var Port, TarPath string
	flag.StringVar(&Port, "port", "./", "listening port")
	flag.StringVar(&TarPath, "path", "./", "absolute file path")
	Port = ":" + Port
	Port = ":5500"
	ln, err := net.Listen("tcp", Port)
	if err != nil {
		fmt.Printf("false listening port: %s", err)
		return
	}

	msghandler.RegisterHandler()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error: %s", err)
			continue
		}

		fmt.Println("new connection")
		go HandleConnection(conn)
	}
}

func HandleConnection(conn net.Conn) {
	header := make([]byte, 4)
	for {
		_, err := conn.Read(header)
		if err != nil {
			if err == syscall.EINVAL {
				continue
			} else if err == io.EOF {
				fmt.Println("client close connection")
				return
			}
			panic(err)
		}

		msgLen, header1 := protocol.Decode_uint16(header)
		msgId, _ := protocol.Decode_uint16(header1)
		var content []byte
		if !(msgLen >= 0 && msgLen < 65535) {
			fmt.Printf("msg len error: %d, connection colse", msgLen)
			conn.Close()
			return
		}

		content = make([]byte, msgLen)
		_, err = conn.Read(content)
		if err != nil {
			if err == syscall.EINVAL {
				continue
			} else if err == io.EOF {
				fmt.Println("client close connection")
				return
			}
			panic(err)
		}
		fmt.Printf("handling msg: %d, len: %d\n", msgId, msgLen)
		msghandler.HandleMsg(msgId, content)
	}
}
