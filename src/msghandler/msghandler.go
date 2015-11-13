package msghandler

import (
	"crypto/md5"
	"fmt"
	"gloger"
	"io"
	"net"
	"os"
	"path/filepath"
	"protocol"
	"runtime"
	"strings"
)

var NeedSrcFiles = make([]string, 0, 5)
var ModTarFiles = make([]string, 0, 5)
var CurFileIdx int
var CurFileCache []byte
var ClientMD5 = md5.New()
var ServMD5 = md5.New()

var TotalRecvPathNum uint32 = 0

const MAX_READ_SIZE = 1024

func ReqFileInfo(Path string, Bpos uint32, conn net.Conn) {
	obj := protocol.CreateServNeedClientData()
	obj.Path = Path
	obj.Bpos = Bpos
	buff := Marshal(uint16(S2C_NEED_FILE_INFO), obj)

	//gloger.GetLoger().Printf("ReqFileInfo: Path: %s,Bpos: %d, CurFileIdx: %d,len(CurFileCache): %d\n", Path, Bpos, CurFileIdx, len(CurFileCache))

	_, err := conn.Write(buff)
	if err != nil {
		panic(err)
	}
}

func c2s_finfo(msgId uint16, msg []byte, conn net.Conn) {
	TotalRecvPathNum++
	fi := protocol.CreateClientFInfo()
	fi.Unmarshal(msg)
	SrcPath := fi.SrcPath
	TarPath := fi.TarPath
	ModTime := fi.ModTime

	if fi.TotalFileNum == TotalRecvPathNum {
		if len(NeedSrcFiles) == 0 {
			fmt.Println("all files updated")
			conn.Close()
			os.Exit(0)
		}
		CurFileIdx = 0
		//CurFileCache = make([]byte, 0, MAX_READ_SIZE)
		Path := NeedSrcFiles[CurFileIdx]
		ServMD5.Reset()
		ReqFileInfo(Path, 0, conn)
	} else {
		finfo, err := os.Stat(TarPath)
		if err != nil || (uint64(finfo.ModTime().UnixNano()) < ModTime) {
			if err != nil {
				if os.IsNotExist(err) {
					gloger.GetLoger().Printf("c2s_finfo: need to create file: %s\n", TarPath)
				} else {
					panic(err)
				}
			} else {
				gloger.GetLoger().Printf("c2s_finfo: need to cover: %s\n", TarPath)
			}
			NeedSrcFiles = append(NeedSrcFiles, SrcPath)
			ModTarFiles = append(ModTarFiles, TarPath)
		}
	}
}

func c2s_file_info_buff(msgId uint16, msg []byte, conn net.Conn) {
	dataBuff := protocol.CreateClientDataBuff()
	dataBuff.Unmarshal(msg)
	CurPath := NeedSrcFiles[CurFileIdx]
	if CurPath != dataBuff.Path {
		panic(fmt.Sprintf("CurPath: %s, sending path: %s", CurPath, dataBuff.Path))
	}
	//CurFileCache = append(CurFileCache, dataBuff.Buff...)
	ServMD5.Write(dataBuff.Buff)
	TarPath := ModTarFiles[CurFileIdx]
	TmpModFile := TarPath + ".tmp"

	Dir := filepath.Dir(TarPath)
	err := os.MkdirAll(Dir, 0666)
	if err != nil {
		panic(err)
	}
	fh, err1 := os.OpenFile(TmpModFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err1 != nil {
		panic(err1)
	}
	defer fh.Close()

	//gloger.GetLoger().Printf("c2s_file_info_buff: writing file: %s, len: %d\n", TmpModFile, len(dataBuff.Buff))

	fh.Write(dataBuff.Buff)

	ReqFileInfo(dataBuff.Path, dataBuff.Bpos, conn)
}

func s2c_need_file_info(msgId uint16, msg []byte, conn net.Conn) {
	obj := protocol.CreateServNeedClientData()
	obj.Unmarshal(msg)

	Path := obj.Path
	Bpos := obj.Bpos
	if Bpos < 0 {
		panic(Bpos)
	}
	if runtime.GOOS == "windows" {
		Path = strings.Replace(Path, "/", "\\", -1)
	}
	fh, err := os.Open(Path)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	//send pieces of file data to server
	finfo, err1 := fh.Stat()
	if err1 != nil {
		panic(err1)
	}

	//gloger.GetLoger().Printf("s2c_need_file_info: Path: %s,Bpos: %d, file size: %d\n", Path, Bpos, finfo.Size())

	if Bpos > uint32(finfo.Size()) {
		panic(Bpos)
	}
	if Bpos == uint32(finfo.Size()) {
		//send file md5sum to server, to check if file completely transferred
		mobj := protocol.CreateFileMD5Info()
		mobj.Path = Path
		mobj.Sum = fmt.Sprintf("%x", ClientMD5.Sum(nil))

		ClientMD5.Reset()
		buff := Marshal(uint16(C2S_FILE_MD5_INFO), mobj)
		_, err := conn.Write(buff)
		if err != nil {
			panic(err)
		}

		gloger.GetLoger().Printf("complete sending file: %s\n", Path)
	} else {
		if Bpos == 0 {
			ClientMD5.Reset()
		}
		//send MAX_READ_SIZE bytes of data of the file to server
		buff := make([]byte, MAX_READ_SIZE)
		n, err := fh.ReadAt(buff, int64(Bpos))
		if err != nil && err != io.EOF {
			panic(err)
		}
		buff = buff[:n]
		dataBuff := protocol.CreateClientDataBuff()
		dataBuff.MsgId = uint16(C2S_FILE_INFO_BUFF)
		dataBuff.Buff = make([]uint8, len(buff))
		copy(dataBuff.Buff, []uint8(buff))
		dataBuff.Bpos = obj.Bpos + uint32(len(buff))
		dataBuff.Total = uint32(finfo.Size())
		dataBuff.Path = obj.Path
		ClientMD5.Write(dataBuff.Buff)
		buff = Marshal(uint16(C2S_FILE_INFO_BUFF), dataBuff)
		_, err = conn.Write(buff)
		if err != nil {
			panic(err)
		}
	}
}

func c2s_file_md5_info(msgId uint16, msg []byte, conn net.Conn) {
	mobj := protocol.CreateFileMD5Info()
	mobj.Unmarshal(msg)
	serv_md5val := fmt.Sprintf("%x", ServMD5.Sum(nil))
	client_md5val := mobj.Sum
	ServMD5.Reset()
	if serv_md5val != client_md5val {
		panic(fmt.Sprintf("file: %s md5 failed: serv_md5val: %s, client_md5val: %s", mobj.Path, serv_md5val, client_md5val))
	}

	//delete old file, rename .tmp file
	TarFile := ModTarFiles[CurFileIdx]
	var err = os.Remove(TarFile)
	if err != nil {
		if os.IsExist(err) {
			panic(err)
		}
	}

	TarTmpFile := TarFile + ".tmp"
	err = os.Rename(TarTmpFile, TarFile)
	if err != nil {
		panic(err)
	}

	//write modified file to path
	CurFileIdx++
	if CurFileIdx < len(NeedSrcFiles) {
		//CurFileCache = make([]byte, 0, MAX_READ_SIZE)
		Path := NeedSrcFiles[CurFileIdx]

		ReqFileInfo(Path, 0, conn)
	} else {
		fmt.Println("all files updated")
		conn.Close()
		os.Exit(0)
	}
}
