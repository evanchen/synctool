package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

func main() {
	servAddr := fmt.Sprintf("%s:%s", ServIp, ServPort)
	conn, err := net.Dial("tcp", servAddr)
	if err != nil {
		log.Printf("failed to connect to server: %s\n", servAddr)
		//return
	}
	defer conn.Close()

	fmt.Printf("connected to server: %s!\n", servAddr)

	cmd := exec.Command("svn", "status", SvnPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to run cmd, reason: %v\n", err)
		return
	}

	fmt.Printf("cmd output: \n%s\n", out.String())

	fmt.Println("extracting files...")
	//var TarFiles
	for {
		line, err := out.ReadString('\n')
		if err != nil {
			break
		}

		line = strings.Replace(line, "\\", "/", -1)
		if strings.HasPrefix(line, "X") {
			fmt.Printf("ignore %s", line)
			continue
		}
		if strings.HasPrefix(line, "?") {
			if strings.Contains(line, ".") {
				fmt.Printf("extracting new file: %s", line)
			} else { //is a folder
				fmt.Printf("ignore %s", line)
			}
		} else if strings.HasPrefix(line, "M") {
			fmt.Printf("extracting modified file: %s", line)
		}

		if strings.
	}
}
