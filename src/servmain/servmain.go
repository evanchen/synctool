package main

import (
	"flag"
	"fmt"
	"gloger"
	"io"
	"msghandler"
	"net"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func main() {
	var Port string
	flag.StringVar(&Port, "port", "./", "listening port")
	flag.Parse()
	Port = ":" + Port
	ln, err := net.Listen("tcp", Port)
	if err != nil {
		fmt.Printf("false listening port: %s", err)
		return
	}
	t := time.Now()
	logName := fmt.Sprintf("serv_%04d%02d%02d_%02d%02d%02d.log", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	gloger.CreateFL(logName)
	msghandler.RegisterHandler()

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("accept error: %s", err)
		return
	}
	defer conn.Close()

	wg.Add(1)
	go HandleConnection(conn)

	wg.Wait()
}

func HandleConnection(conn net.Conn) {
	defer wg.Done()

	for {
		msgId, content, err := msghandler.DoRecv(conn)

		if err != nil {
			if err == syscall.EINVAL {
				continue
			} else if err == io.EOF {
				fmt.Println("connection closed")
				conn.Close()
				return
			}
			fmt.Println(err)
			return
		} else {
			//fmt.Printf("handling msg: %d, len: %d\n", msgId, len(content))
		}

		msghandler.HandleMsg(msgId, content, conn)
	}
}
