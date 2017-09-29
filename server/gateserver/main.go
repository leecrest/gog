// 网关服务器

package main

import (
	"flag"
	"github.com/leecrest/gog/engine/config"
	"os"
	"fmt"
)

type GateCfg struct {
	NetID byte
	ServerID uint32			// 服务器组编号
	WorkSpace string		// 服务器组的工作路径
	LogPath string 			// 日志文件路径，与工作路径的相对路径
	RemoteAddr string		// 对外接口
	LocalAddr string		// 对内接口
}

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
	return true
}


// 启动参数：-cfg配置文件路径 -id网关编号
func main() {
	// 解析命令行参数
	cfg := flag.String("cfg", "xgame.cfg", "config file")
	id := flag.Int("id", 0, "gate id")
	flag.Parse()

	g_config.NetID = byte(*id)

	var ret = loadConfig(*cfg)
	if !ret {
		return
	}

	os.Chdir(g_config.WorkSpace)

	// 启动对外监听socket
	fmt.Println("启动对外接口")
	ret = initRemote(g_config.RemoteAddr)
	if !ret {
		fmt.Println("对外接口启动失败")
		return
	}

	// 启动进程间通信的socket
	fmt.Println("启动对内接口")
	ret = initLocal(g_config.LocalAddr)
	if !ret {
		fmt.Println("对内接口启动失败")
		return
	}
}