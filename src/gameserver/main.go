package main

import (
	"flag"
	"github.com/leecrest/gog/src/gameserver/gameserver"
)


// 启动参数：-cfg=配置文件路径 -id=进程编号
func main() {
	// 解析命令行参数
	path := flag.String("cfg", "server.cfg", "config file")
	id := flag.Int("id", 0, "gate id")
	flag.Parse()

	gameserver.Run(*path, byte(*id))
}