package gate

/*
说明：
1、主线程用于本地rpc的监听
2、启动1个协程监听外部端口
3、每一个外部客户端链接，启动1个读协程
4、每一个本地链接，启动1个读协程
 */

import (
	"os"
	"fmt"
	"github.com/leecrest/gog/src/engine/config"
)

var g_config GateCfg

func loadConfig(path string) bool {
	var cfg, err = config.INILoad(path)
	if err != nil {
		logError("load config error: " + err.Error())
		return false
	}
	var section = "default"
	g_config.serverID = cfg.ReadUint32(section, "server", 0)
	g_config.workSpace = cfg.Read(section, "workspace", "./")
	g_config.logPath = cfg.Read(section, "log", "./log")
	g_config.console = cfg.ReadInt(section, "console", 0) == 1

	section = "gate"
	g_config.num = (byte)(cfg.ReadUint32(section, "num", 1))
	g_config.ip = cfg.Read(section, "ip", "ip")
	g_config.wanPort = cfg.ReadUint32(section, "wan_port", 0)
	g_config.lanPort = cfg.ReadUint32(section, "lan_port", 0)

	var port = g_config.wanPort + uint32(g_config.id)
	g_config.wanAddr = fmt.Sprintf("%s:%d", g_config.ip, port)
	port = g_config.lanPort + uint32(g_config.id)
	g_config.lanAddr = fmt.Sprintf("%s:%d", g_config.ip, port)
	return true
}


// 启动参数：-cfg=配置文件路径 -id=网关编号
func Run(path string, id byte) {
	g_config.id = id
	var err = loadConfig(path)
	if !err {
		return
	}
	os.Chdir(g_config.workSpace)
	logInit()

	// 启动对外监听socket
	err = initClient(g_config.wanAddr, false)
	if !err {
		return
	}

	// 启动进程间通信
	err = initGame(g_config.lanAddr, true)
	if !err {
		return
	}
}