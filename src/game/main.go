package main

import (
	"flag"
	"github.com/leecrest/gog/src/game/game"
)


// 启动参数：-cfg=配置文件路径 -id=进程编号 -name=进程名称
func main() {
	// 解析命令行参数
	path := flag.String("cfg", "server.cfg", "config file")
	id := flag.Int("id", 0, "gate id")
	name := flag.String("name", "game", "server name")
	flag.Parse()

	game.Run(byte(*id), *name, *path)
}