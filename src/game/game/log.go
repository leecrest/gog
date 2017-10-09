package game

import (
	"os"
	"fmt"
	"log"
)

var g_logFile *os.File
var g_logger *log.Logger
var g_logPrefix string

// 初始化日志文件
func logInit() (bool) {
	var err = os.MkdirAll(g_config.logPath, os.ModeDir)
	if err != nil {
		panic(err)
		return false
	}
	var file = g_config.logPath + "/game.log"
	g_logFile, err = os.OpenFile(file, os.O_APPEND | os.O_CREATE, os.ModeAppend)
	if err != nil {
		panic(err)
		return false
	}
	g_logPrefix = fmt.Sprintf("(%d)", g_config.id)
	g_logger = log.New(g_logFile, "", log.LstdFlags)
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
	var msg = g_logPrefix + " [INFO] " + fmt.Sprintf(format, a...)
	g_logger.Println(msg)
	if g_config.console {
		fmt.Println(msg)
	}
}

func logError(format string, a ...interface{}) {
	var msg = g_logPrefix + " [ERROR] " + fmt.Sprintf(format, a...)
	g_logger.Println(msg)
	if g_config.console {
		fmt.Println(msg)
	}
}