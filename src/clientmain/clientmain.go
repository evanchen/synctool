package main

import (
	"flag"
	"fmt"
	//"net"
	"pathanalysis"
	"protocol"
)

func main() {
	var ServIp, ServPort, TarPath, IgnoreList, IncludeList string
	flag.StringVar(&ServIp, "ip", "./", "server ip")
	flag.StringVar(&ServPort, "port", "./", "server port")
	flag.StringVar(&TarPath, "path", "./", "absolute file path")
	flag.StringVar(&IgnoreList, "ignore-dir", "./", "ignore directory: setting;.svn;common")
	flag.StringVar(&IncludeList, "include", "./", "include surfix files: lua;cpp")
	flag.Parse()

	fmt.Println(ServIp, ServPort, TarPath, IgnoreList, IncludeList)
	/*
		servAddr := fmt.Sprintf("%s:%s", ServIp, ServPort)
		conn, err := net.Dial("tcp", servAddr)
		if err != nil {
			fmt.Printf("failed to connect to server: %s\n", servAddr)
			return
		}
		defer conn.Close()

		fmt.Printf("connected to server: %s!\n", servAddr)
	*/
	ServIp, ServPort, TarPath, IgnoreList, IncludeList = "192.168.1.98", "5500", "e:\\trunk\\logic", ".svn;setting;common", "*.lua;*.cpp"
	listFInfo := pathanalysis.DoAnalysis(TarPath, IgnoreList, IncludeList)
	pl := protocol.CreateFInfoList()
	for i := 0; i < len(listFInfo); i++ {
		fi := protocol.CreateFInfo()
		fi.Path = listFInfo[i].Path
		fi.Modtime = uint64(listFInfo[i].ModifyTime)
		pl.FinfoList = append(pl.FinfoList, *fi)
	}

	pl.Marshal()
	fmt.Println("done!")
}
