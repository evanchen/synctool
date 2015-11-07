package pathanalysis

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var g_logger *log.Logger

//create a log file and log.Logger
func createFL(fname string) (*os.File, *log.Logger) {
	path := "logs/" + fname
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("failed to create logfile: %s: %s", fname, err.Error())
		return nil, nil
	}
	l := log.New(f, "", log.LstdFlags)

	return f, l
}

// ignore directory list
var g_IgnoreList map[string]bool = make(map[string]bool)

// include file surfixes list
var g_IncludeList map[string]bool = make(map[string]bool)

type FInfo struct {
	Name       string
	Path       string
	Type       byte
	ModifyTime int64
}

var g_Ch = make(chan *FInfo, 5)
var g_RootName string

func Init(path, ignoreList, includeList string) {
	_, lg := createFL("output.log")
	if lg == nil {
		panic("logger failed")
	}
	g_logger = lg
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
func DoAnalysis(path, ignoreList, includeList string) (List []*FInfo) {
	Init(path, ignoreList, includeList)
	fmt.Println(path, ignoreList, includeList, g_RootName, g_IgnoreList, g_IncludeList)
	go func() {
		defer close(g_Ch)
		filepath.Walk(path, walkFunc)
	}()

	for v := range g_Ch {
		List = append(List, v)
		g_logger.Printf("%s %s %d %d\n", v.Path, v.Name, v.Type, v.ModifyTime)
	}

	return List
}

//walk function
func walkFunc(path string, info os.FileInfo, err error) error {
	if err != nil {
		panic(err)
	}
	// get root directory: for example, "path: /home/chenwenqiang/logic/module/act/act_pto.lua" -> "logic/module/act/act_pto.lua",
	// assuming root name is "logic"

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

	fi := new(FInfo)
	fi.Name = name
	fi.ModifyTime = info.ModTime().UnixNano()
	fi.Path = path
	if isdir {
		fi.Type = 1
	} else {
		fi.Type = 2
	}

	g_Ch <- fi

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
		//g_logger.Printf("ignore file: %s\n", name)
		return true
	}

	return false
}

func FilterDir(path, name string) bool {
	if name == "." || name == ".." {
		//g_logger.Printf("ignore diretory: %s, %s\n", path, name)
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
			//g_logger.Printf("ignore diretory: %s, %s\n", path, name)
			return true
		}
	}

	return false
}
