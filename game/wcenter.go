package game

import (
	"c-server/blockchain/eth"
	"c-server/cache"
	"c-server/model"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"reflect"
	"strconv"
)

//用户中心

type Wcenter struct {
	Name string
}

//type WebcenterActions interface {
//	Login(interface{}) []byte
//	Register(interface{}) interface{}
//}

func (wb Wcenter) Info(req ClientMessage) {
	res := ResponseData{}
	res.Data = make(map[string]interface{})
	loginfo, ok := req.Data["msg"].(map[string]interface{})
	/*
			req.response({
		cake: toCake(cake),
		balance,
		credit: sqlresult[0].credit,
		invite: deepGetNumber(await redis.get(INV)),
		inviteGrow: conditions[lvl],
		verified,//0未认证,1已认证,-1已提交申请
		hasAirdropCake: Number(hasAirdropCake)
		});
	*/
	if !ok {
		req.error("错误")
	}
	token := loginfo["token"]
	p, e := model.GetPlayerById(1)
	if e != nil {
		req.error("玩家不存在")
		return
	}
	cake := p.GetCoin("coin1")
	fmt.Println(token)
	res.Put("cake", cake)
	res.Put("balance", p.GetNeoBalance())
	res.Put("hasAirdropCake", 0)
	verified, err := cache.RedisClient.Get("player:" + strconv.Itoa(1) + ":verified").Result()
	if err != nil {
		verified = "0"
	}
	account := common.HexToAddress("0x87c5985dBc32650F20E275453DddE97e5D44eC82")
	balance, err := eth.Client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance) // 25893180161173005034
	res.Put("verified", verified)
	req.response(res)
}
func (wb Wcenter) ChargeAndWithdraw(req ClientMessage) {
	res := ResponseData{}
	res.Data = make(map[string]interface{})
	/*
		req.response({
						balance,
						nnc: await player.getToken("nnc"),
						tax: deepGetNumber(GameData.getGlobalDataValue("tax_withdraw_eth")),
						chargeRecords,
						withdrawRecords,
						deniedWithdraws
					});
	*/
	//res.Put("balance")
	res.Put("chargeRecords", []interface{}{})
	res.Put("withdrawRecords", []interface{}{})
	req.response(res)
}

func (wb *Wcenter) RegisterController() map[string]reflect.Value {
	cont := make(map[string]reflect.Value)
	v := reflect.ValueOf(wb).Elem()
	t := reflect.TypeOf(wb).Elem()
	for i := 0; i < v.NumMethod(); i++ {
		action := v.Method(i)
		actionName := t.Method(i).Name
		cont[actionName] = action
	}
	return cont
}
