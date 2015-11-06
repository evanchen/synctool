package main

import (
	"fmt"
	"net"
	"pathanalysis"
	"protocol"
)

func main() {
	var ServIp, ServPort, TarPath, IgnoreList, IncludeList string
	flag.StringVar(&ServIp, "ip", "./", "server ip")
	flag.StringVar(&ServPort, "port", "./", "server port")
	flag.StringVar(&TarPath, "path", "./", "absolute file path")
	flag.StringVar(&IgnoreList, "ignore", "./", "ignore surfix files: lua;cpp;hpp")
	flag.StringVar(&IncludeList, "include", "./", "include surfix files: lua;cpp")
	flag.Parse()

	servAddr := fmt.Sprintf("%s:%s", ServIp, ServPort)
	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		fmt.Printf("failed to connect to server: %s\n", servAddr)
		return
	}
	defer conn.Close()

	fmt.Printf("connected to server: %s!\n", servAddr)

	listFInfo := pathanalysis.DoAnalysis(TarPath, IgnoreList, IncludeList)
	pl := protocol.CreateFInfoList()
	for i := 0; i < len(listFInfo); i++ {
		fi := protocol.CreateFInfo()
		fi.Path = listFInfo[i].Path
		fi.Modtime = listFInfo[i].ModifyTime
		pl.FinfoList = append(pl.FinfoList, *fi)
	}

	buff := pl.Marshal()
}
