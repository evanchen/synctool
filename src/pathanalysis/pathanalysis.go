package pathanalysis

import (
	"path/filepath"
	"os"
	"strings"
)

// ignore surfixes list, only for files but not directories
var g_IgnoreList map[string]bool = make(map[string]bool)
// include surfixes list, only for files but not directories
var g_IncludeList map[string]bool = make(map[string]bool)

type FInfo struct{
	Name string
	Path string 
	Type byte
	ModifyTime int64
}

var g_Ch = make(chan *FInfo,5)
var g_RootName string

func Init(path,ignoreList,includeList string) {
	name := filepath.Base(path)
	g_RootName = name

	f := func extract(tolist map[string]bool,strList) {
		arrstr := strings.Split(strList, ";")
		for i := 0; i < len(arrstr); i ++ {
			lidx := strings.LastIndex(arrstr[i], ".")  // *.lua -> lua
			if lidx >= 0 {
				arrstr[i] = string(arrstr[lidx:])
			}
			arrstr[i] = strings.TrimSpace(arrstr[i])
			tolist[ arrstr[i] ] = true
		}
	}
	f(g_IgnoreList,ignoreList)
	f(g_IncludeList,includeList)
}

//walk through child directories and files
func DoAnalysis(path,ignoreList,includeList string) (List []*FInfo){
	Init(path,ignoreList,includeList)
	go func(){
		defer close(g_Ch)
		filepath.Walk(path,walkFunc)
	}()

	for v := range g_Ch {
		List = appand(List,v)
	}

	return List
}


//walk function
func walkFunc(path string, info os.FileInfo, err error) error
	if err != nil {
		panic(err)
	}
	// get root directory: for example, "path: /home/chenwenqiang/logic/module/act/act_pto.lua" -> "logic/module/act/act_pto.lua", 
	// assuming root name is "logic"
	baseName := filepath.Base(path)

	if _,ok := g_IgnoreList[baseName]; ok {
		return 
	}

	if _,ok := g_IncludeList[baseName]; !ok {
		return 
	}

	idx := strings.Index(path, g_RootName)
	if idx < 0 {
		panic("illegal path")
	}
	path = path[idx:]

	fi := new(FInfo)
	fi.Name = info.Name()
	fi.ModifyTime = info.ModTime().UnixNano()
	fi.Path = path
	if info.IsDir() {
		fi.Type = 1
	} else {
		fi.Type = 2
	}

	g_Ch <- fi

	return err
}

