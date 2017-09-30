package main

import (
	"github.com/leecrest/gog/src/gateserver/gateserver"
	"flag"
)



// 启动参数：-cfg=配置文件路径 -id=网关编号
func main() {
	// 解析命令行参数
	path := flag.String("cfg", "server.cfg", "config file")
	id := flag.Int("id", 0, "gate id")
	flag.Parse()

	gateserver.Run(*path, byte(*id))
}