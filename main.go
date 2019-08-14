package main

import (
	"c-server/blockchain/eth"
	"c-server/conf"
	"c-server/game"
	"c-server/server"
	"math/big"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	eth.Init()
	game.Init()
	eth.Blo = make(chan *big.Int, 100000)
	eth.Save = make(chan *big.Int, 100000)
	//go eth.WriteToRedis()
	//go eth.SyncBlock()
	go eth.SyncTx()
	go eth.StartSync()

	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")
}
