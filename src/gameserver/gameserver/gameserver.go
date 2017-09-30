package gameserver


import (
	"os"
	"github.com/leecrest/gog/src/engine/config"
	"fmt"
)

var g_config GameCfg

func loadConfig(path string) (bool) {
	var cfg, err = config.INILoad(path)
	if err != nil {
		panic(err)
		return false
	}
	g_config.ServerID = cfg.ReadUint32("default", "server", 0)
	g_config.WorkSpace = cfg.Read("default", "workspace", "./")
	g_config.LogPath = cfg.Read("default", "log", "./log")
	g_config.LocalAddr = cfg.Read("gateserver", "local", "")
	g_config.LogPrint = cfg.ReadInt("default", "print", 0) == 1
	return true
}



// 启动参数：-cfg=配置文件路径 -id=进程编号
func Run(path string, id byte) {
	g_config.GameID = id
	g_config.GameIDStr = fmt.Sprintf(" (%d) ", g_config.GameID)
	var ret = loadConfig(path)
	if !ret {
		return
	}
	os.Chdir(g_config.WorkSpace)
	logInit()
}