package pathanalysis

import (
	"fmt"
	"gloger"
	"msghandler"
	"net"
	"os"
	"path/filepath"
	"protocol"
	"runtime"
	"strings"
)

// ignore directory list
var g_IgnoreList map[string]bool = make(map[string]bool)

// include file surfixes list
var g_IncludeList map[string]bool = make(map[string]bool)

var g_Ch = make(chan *protocol.ClientFInfo, 5)
var g_RootName string

func Init(path, ignoreList, includeList string) {

	name := filepath.Base(path)
	g_RootName = name

	ParseBlockDirList(g_IgnoreList, ignoreList)
	ParseBlockSurfixList(g_IncludeList, includeList)
}

func ParseBlockDirList(tolist map[string]bool, strList string) {
	arrstr := strings.Split(strList, ";")
	for i := 0; i < len(arrstr); i++ {
		arrstr[i] = strings.TrimSpace(arrstr[i])
		tolist[arrstr[i]] = true
	}
}

func ParseBlockSurfixList(tolist map[string]bool, strList string) {
	arrstr := strings.Split(strList, ";")
	for i := 0; i < len(arrstr); i++ {
		lidx := strings.LastIndex(arrstr[i], ".") // *.lua -> .lua
		if lidx >= 0 {
			arrstr[i] = arrstr[i][lidx:]
		}
		arrstr[i] = strings.TrimSpace(arrstr[i])
		tolist[arrstr[i]] = true
	}
}

//walk through child directories and files
func DoAnalysis(path, ignoreList, includeList, ServPath string, conn net.Conn) {
	Init(path, ignoreList, includeList)
	fmt.Println(path, ignoreList, includeList, g_RootName, g_IgnoreList, g_IncludeList)
	go func() {
		defer close(g_Ch)
		filepath.Walk(path, walkFunc)
	}()

	list := make([]*protocol.ClientFInfo, 0, 8)
	//pack informations
	for v := range g_Ch {
		//gloger.GetLoger().Printf("%s %d\n", v.SrcPath, v.ModTime)

		v.TarPath = filepath.Join(ServPath, v.TarPath)

		if runtime.GOOS == "windows" {
			v.TarPath = strings.Replace(v.TarPath, "\\", "/", -1)
		}

		list = append(list, v)
	}

	size := len(list)
	if size == 0 {
		fmt.Printf("Path: %s no files\n", path)
		conn.Close()
		os.Exit(0)
	} else {
		for i := 0; i < size; i++ {
			fi := list[i]
			fi.TotalFileNum = uint32(size)
			buff := msghandler.Marshal(uint16(msghandler.C2S_FINFO), fi)
			_, err := conn.Write(buff)
			if err != nil {
				panic(fmt.Sprintf("connection write error: %s", err))
			}
		}
	}
}

//walk function
func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	// get root directory: for example, "path: /home/chenwenqiang/logic/module/act/act_pto.lua" -> "logic/module/act/act_pto.lua",
	// assuming root name is "logic"
	orgPath := path
	isdir := info.IsDir()
	name := info.Name()

	idx := strings.Index(path, g_RootName)
	if idx < 0 {
		panic("illegal path")
	}
	path = path[idx:]

	if isdir {
		if FilterDir(path, name) {
			return err
		}
	} else {
		if FilterFileSurfix(path, name) {
			return err
		}
	}

	if !isdir {
		if info.Size() == 0 {
			gloger.GetLoger().Printf("ignore empty file : %s\n", orgPath)
			return err
		}
		fi := protocol.CreateClientFInfo()
		fi.ModTime = uint64(info.ModTime().UnixNano())
		fi.TarPath = path
		fi.SrcPath = orgPath
		g_Ch <- fi
	}

	return err
}

func FilterFileSurfix(path, name string) bool {
	dir := filepath.Dir(path)
	if FilterDir(dir, name) {
		return true
	}
	lidx := strings.LastIndex(name, ".") // *.lua -> .lua
	var surfix string
	if lidx >= 0 {
		surfix = name[lidx:]
	}
	if _, ok := g_IncludeList[surfix]; !ok {
		return true
	}

	return false
}

func FilterDir(path, name string) bool {
	if name == "." || name == ".." {
		return true
	}

	var arrpath []string
	if runtime.GOOS == "windows" {
		arrpath = strings.Split(path, "\\")
	} else {
		arrpath = strings.Split(path, "/")
	}
	for i := 0; i < len(arrpath); i++ {
		if _, ok := g_IgnoreList[arrpath[i]]; ok {
			return true
		}
	}

	return false
}
