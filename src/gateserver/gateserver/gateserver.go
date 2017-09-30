package gateserver

import (
	"github.com/leecrest/gog/src/engine/config"
	"os"
)

var g_config GateCfg

func loadConfig(path string) (bool) {
	var cfg, err = config.INILoad(path)
	if err != nil {
		panic(err)
		return false
	}
	g_config.ServerID = cfg.ReadUint32("default", "server", 0)
	g_config.WorkSpace = cfg.Read("default", "workspace", "./")
	g_config.LogPath = cfg.Read("default", "log", "./log")
	g_config.RemoteAddr = cfg.Read("gateserver", "remote", "")
	g_config.LocalAddr = cfg.Read("gateserver", "local", "")
	g_config.LogPrint = cfg.ReadInt("default", "print", 0) == 1
	return true
}


// 启动参数：-cfg=配置文件路径 -id=网关编号
func Run(path string, id byte) {
	g_config.NetID = id
	var ret = loadConfig(path)
	if !ret {
		return
	}
	os.Chdir(g_config.WorkSpace)
	logInit()

	// 启动对外监听socket
	ret = initRemote(g_config.RemoteAddr)
	if !ret {
		logInfo("init remote error")
		return
	}

	// 启动进程间通信的socket
	ret = initLocal(g_config.LocalAddr)
	if !ret {
		logInfo("init local error")
		return
	}
}