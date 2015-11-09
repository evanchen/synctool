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
)

var wg sync.WaitGroup

func main() {
	var ServIp, ServPort, TarPath, ServPath, IgnoreList, IncludeList string
	flag.StringVar(&ServIp, "ip", "./", "server ip")
	flag.StringVar(&ServPort, "port", "./", "server port")
	flag.StringVar(&TarPath, "cpath", "./", "client absolute file path")
	flag.StringVar(&ServPath, "spath", "./", "server absolute file path")
	flag.StringVar(&IgnoreList, "ignore-dir", "./", "ignore directory: setting;.svn;common")
	flag.StringVar(&IncludeList, "include", "./", "include surfix files: lua;cpp")
	flag.Parse()

	ServIp, ServPort, TarPath, ServPath, IgnoreList, IncludeList = "127.0.0.1", "5500", "e:\\trunk\\logic", "f:\\test", ".svn;setting;common", "*.lua;*.cp"
	fmt.Println(ServIp, ServPort, TarPath, ServPath, IgnoreList, IncludeList)

	gloger.CreateFL("logclient.log")
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

		fmt.Printf("handling msg: %d, len: %d\n", msgId, len(content))
		msghandler.HandleMsg(msgId, content, conn)
	}
}
