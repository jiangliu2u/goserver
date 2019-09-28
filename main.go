package main

import (
	"c-server/blockchain/eth"
	"c-server/conf"
	"c-server/game"
	"c-server/server"
	"sync"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	eth.Init()
	game.Init()
	//go eth.WriteToRedis()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	//go eth.SyncBlock()
	go eth.StartSync(wg)
	wg.Wait()

	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
