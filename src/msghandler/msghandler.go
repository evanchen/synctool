package msghandler

import (
	"gloger"
	"net"
	"os"
	"protocol"
)

func c2s_finfo(msgId uint16, msg []byte, conn net.Conn) {
	infoList := protocol.CreateFInfoList()
	infoList.Unmarshal(msg)
	for i := 0; i < len(infoList.FinfoList); i++ {
		Path := infoList.FinfoList[i].Path
		ModTime := infoList.FinfoList[i].ModTime

		// if file not exsit, create one
		finfo, err := os.Stat(Path)
		if err != nil || (uint64(finfo.ModTime().UnixNano()) < ModTime) {
			if err != nil {
				gloger.GetLoger().Printf("c2s_finfo: error: %v\n", err)
			}
			gloger.GetLoger().Printf("c2s_finfo: need to cover: %s\n", Path)

			obj := protocol.CreateFilePath()
			obj.Path = Path
			buff := Marshal(uint16(S2C_UPDATE_FILE), obj)
			_, err = conn.Write(buff)
			if err != nil {
				panic(err)
			}
		}
	}
}

func s2c_finfo(msgId uint16, msg []byte, conn net.Conn) {

}

func c2s_update_file(msgId uint16, msg []byte, conn net.Conn) {

}

func s2c_update_file(msgId uint16, msg []byte, conn net.Conn) {
	obj := protocol.CreateFilePath()
	obj.Unmarshal(msg)

	//send modified file to server
	gloger.GetLoger().Printf("send file: %s\n", obj.Path)
}

func s2c_done(msgId uint16, msg []byte, conn net.Conn) {
	conn.Close()
}
