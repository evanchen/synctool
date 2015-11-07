package main

import (
	"coding"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type MsgHandlerFunc func(uint16, []byte)

var g_msgHandlers = make(map[uint16]MsgHandlerFunc)

const (
	Port = ":1970"
)

func main() {
	ln, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("false listening port: %s", err)
		return
	}

	RegisterMsgHandler()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accept error: %s", err)
			continue
		}

		go HandleConnection(conn)
	}
}

func RegisterMsgHandler() {

}

func HandleConnection(conn net.Conn) {
	header := make([]byte, 4)
	for {
		rz, err := conn.Read(header)
		if err != nil {
			if err == syscall.EINVAL {
				continue
			} else if err == io.EOF {
				fmt.Println("client close connection")
				return
			}
			panic(err)
		}

		msgId, header1 := coding.decode_uint16(header)
		msgLen, _ := coding.decode_uint16(header1)
		var content []byte

		f := GetMsgHandler(msgId)
		if !f {
			fmt.Printf("no msg handler for : %d \n", msgId)
			if msgLen > 0 {
				content = make([]byte, msgLen)
				conn.Read(content)
				continue
			}
		} else if msgLen > 0 {
			content = make([]byte, msgLen)
			conn.Read(content)
		}

		f(msgId, content)
	}
}

func GetMsgHandler(msgId uint16) MsgHandlerFunc {
	f, ok := g_msgHandlers[msgId]
	if !ok {
		return nil
	}

	return f
}
