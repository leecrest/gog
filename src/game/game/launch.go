package game

/*
说明：
1、
*/


import (
	"os"
	"github.com/leecrest/gog/src/engine/config"
	"fmt"
	"time"
)

var g_config GameCfg

func loadConfig(path string) (bool) {
	var cfg, err = config.INILoad(path)
	if err != nil {
		panic(err)
		return false
	}

	var section = "default"
	g_config.serverID = cfg.ReadUint32(section, "server", 0)
	g_config.workSpace = cfg.Read(section, "workspace", "./")
	g_config.logPath = cfg.Read(section, "log", "./log")
	g_config.console = cfg.ReadInt(section, "console", 0) == 1

	section = "gate"
	g_config.gateNum = (byte)(cfg.ReadUint32(section, "num", 1))
	g_config.gateIP = cfg.Read(section, "ip", "")
	g_config.gatePort = cfg.ReadUint32(section, "lan_port", 0)

	section = "game"
	g_config.rpcPort = cfg.ReadUint32(section, "rpc_port", 0)

	section = g_config.name
	g_config.preload = cfg.Read(section, "preload", "preload.go")

	var port uint32 = g_config.rpcPort + uint32(g_config.id)
	g_config.rpcAddr = fmt.Sprintf("%s:%d", g_config.gateIP, port)
	return true
}



// 启动参数：-cfg=配置文件路径 -id=进程编号
func Run(id byte, name string, path string) {
	g_config.id = id
	g_config.name = name
	var ret = loadConfig(path)
	if !ret {
		return
	}
	os.Chdir(g_config.workSpace)
	ret = logInit()
	if !ret {
		return
	}
	ret = gateInit()
	if !ret {
		return
	}
	ret = rpcInit()
	if !ret {
		return
	}

	// 启动game的逻辑层代码
	for {
		time.Sleep(1000)
	}
}