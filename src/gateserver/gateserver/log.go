package gateserver

import (
	"os"
	"fmt"
	"log"
)

var g_logFile *os.File
var g_logger *log.Logger

// 初始化日志文件
func logInit() (bool) {
	var err = os.MkdirAll(g_config.LogPath, os.ModeDir)
	if err != nil {
		panic(err)
		return false
	}
	var file = g_config.LogPath + "/gateserver.log"
	g_logFile, err = os.OpenFile(file, os.O_APPEND | os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
		return false
	}
	var prefix = fmt.Sprintf(" (%d) ", g_config.NetID)
	g_logger = log.New(g_logFile, prefix, log.LstdFlags)
	return true
}

func logClose() {
	if g_logFile == nil {
		return
	}
	g_logFile.Close()
	g_logFile = nil
	g_logger = nil
}

func logInfo(format string, a ...interface{}) {
	var msg = fmt.Sprintf(format, a...)
	g_logger.Println(msg)
	if g_config.LogPrint {
		fmt.Println(msg)
	}
}