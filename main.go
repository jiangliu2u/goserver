package main

import (
	"c-server/blockchain/ethereum"
	"c-server/conf"
	"c-server/game"
	"c-server/server"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	ethereum.Init()
	game.Init()
	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
