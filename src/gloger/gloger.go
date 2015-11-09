package gloger

import (
	"log"
	"os"
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

func GetLoger() *log.Logger {
	return g_logger
}
