package main

import (
	"flag"
	"fmt"
	"msghandler"
	"net"
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

	ServIp, ServPort, TarPath, IgnoreList, IncludeList = "127.0.0.1", "5500", "e:\\trunk\\logic", ".svn;setting;common", "*.lua;*.cp"
	fmt.Println(ServIp, ServPort, TarPath, IgnoreList, IncludeList)

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
		fi.Modtime = uint64(listFInfo[i].ModifyTime)
		pl.FinfoList = append(pl.FinfoList, *fi)
	}

	buff := pl.Marshal()
	csz := len(buff)
	tsz := 4 + csz
	total := make([]byte, tsz)
	tmp := protocol.Encode_uint16(uint16(csz))
	copy(total, tmp)
	fmt.Println(csz, tmp, total[:4])

	msgId := uint16(msghandler.C2S_FINFO)
	tmp = protocol.Encode_uint16(msgId)
	copy(total[2:], tmp)
	fmt.Println(msgId, tmp, total[:4])

	copy(total[4:], buff)
	_, err = conn.Write(total)
	if err != nil {
		fmt.Printf("connection write error: %s", err)
	}
}
