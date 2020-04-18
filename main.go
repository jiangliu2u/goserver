package main

import (
	"c-server/conf"
	"c-server/game"
	"c-server/server"
	"time"
)

func main() {
	// 从配置文件读取配置
	conf.Init()
	// eth.Init()
	game.WSClients = make(map[string]*game.WsConnection)
	game.Init()
	// wg := new(sync.WaitGroup)
	// wg.Add(1)
	// go eth.StartSync(wg)
	// wg.Wait()
	go func() {
		d := time.Duration(time.Second * 2)

		t := time.NewTicker(d)
		defer t.Stop()

		for {
			select {
			case <-t.C:
				for _, v := range game.WSClients {
					data := make(map[string]interface{})
					data["hahaha"] = "jiangliu2u"
					v.PushMessage("S_dddd", data)
				}
			}

		}
	}()
	// 装载路由
	r := server.NewRouter()
	r.Run(":3000")

}
