package main

import (
	"os"
	"time"
	"fmt"
)

var g_log *os.File

// 初始化日志文件
func logInit() (bool) {
	var err = os.MkdirAll(g_config.LogPath, os.ModeDir)
	if err != nil {
		panic(err)
		return false
	}
	var file = g_config.LogPath + "/gateserver.log"
	g_log, err = os.OpenFile(file, os.O_APPEND | os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
		return false
	}
	return true
}

func logClose() {
	if g_log == nil {
		return
	}
	g_log.Close()
	g_log = nil
}

func logInfo(format string, a ...interface{}) {
	var msg = fmt.Sprintf(format, a...)
	var txt = time.Now().Format("2006-01-02 15:04:05")
	txt += g_config.NetStr + msg
	g_log.WriteString(txt + "\n")
	g_log.Sync()
	if g_config.LogPrint {
		fmt.Println(txt)
	}
}