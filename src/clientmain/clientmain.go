package main

import (
	"flag"
	"fmt"
	"gloger"
	"io"
	"msghandler"
	"net"
	"pathanalysis"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

func main() {
	var ServIp, ServPort, TarPath, ServPath, IgnoreList, IncludeList string
	flag.StringVar(&ServIp, "ip", "127.0.0.1", "server ip")
	flag.StringVar(&ServPort, "port", "5500", "server port")
	flag.StringVar(&TarPath, "cpath", "/home/tarpath", "client absolute file path")
	flag.StringVar(&ServPath, "spath", "/home/tmp", "server absolute file path")
	flag.StringVar(&IgnoreList, "ignore-dir", ".svn", "ignore directory")
	flag.StringVar(&IncludeList, "include", "*.lua;*.h", "include surfix files")
	flag.Parse()

	fmt.Println(ServIp, ServPort, TarPath, ServPath, IgnoreList, IncludeList)

	t := time.Now()
	logName := fmt.Sprintf("client_%04d%02d%02d_%02d%02d%02d.log", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	gloger.CreateFL(logName)
	msghandler.RegisterHandler()

	servAddr := fmt.Sprintf("%s:%s", ServIp, ServPort)
	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		fmt.Printf("failed to connect to server: %s\n", servAddr)
		return
	}
	defer conn.Close()

	fmt.Printf("connected to server: %s!\n", servAddr)

	wg.Add(2)
	//analyze file informations and send to server
	go func() {
		defer wg.Done()
		pathanalysis.DoAnalysis(TarPath, IgnoreList, IncludeList, ServPath, conn)
	}()

	//receive msg from server
	go func() {
		defer wg.Done()
		DoRecvMsg(conn)
	}()

	wg.Wait()
}

func DoRecvMsg(conn net.Conn) {
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
		}

		//fmt.Printf("handling msg: %d, len: %d\n", msgId, len(content))
		msghandler.HandleMsg(msgId, content, conn)
	}
}
