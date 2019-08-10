package game

import (
	"fmt"
	"reflect"
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
	fmt.Println(token)

	req.response(res)
}
func (wb Wcenter) ChargeAndWithdraw(req ClientMessage) {
	res := ResponseData{}
	res.Data = make(map[string]interface{})
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
