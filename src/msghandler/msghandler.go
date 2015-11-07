package msghandler

import (
	//"fmt"
	"log"
	"os"
	"protocol"
)

var g_logger *log.Logger

//create a log file and log.Logger
func CreateFL(fname string) {
	path := fname
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("failed to create logfile: %s: %s", fname, err.Error())
		return
	}
	g_logger = log.New(f, "", log.LstdFlags)
}

func c2s_finfo(msgId uint16, msg []byte) {
	infoList := protocol.CreateFInfoList()
	infoList.Unmarshal(msg)
	for i := 0; i < len(infoList.FinfoList); i++ {
		g_logger.Println(infoList.FinfoList[i].Path, infoList.FinfoList[i].Modtime)
	}
}

func s2c_finfo(msgId uint16, msg []byte) {

}

func c2s_update_file(msgId uint16, msg []byte) {

}

func s2c_update_file(msgId uint16, msg []byte) {

}
