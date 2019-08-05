package game

import (
	"c-server/util"
	"encoding/json"
	"fmt"
	"reflect"
)

//用户中心

type Webcenter struct {
	Name string
}

type WebcenterActions interface {
	Login(interface{}) []byte
	Register(interface{}) interface{}
}

func (wb Webcenter) Login(message interface{}) []byte {
	var res ResponseMessage
	var msg ResponseData
	res.Name = "S_response"

	fmt.Println(message)
	decode := message.(map[string]interface{})
	if util.HasAttr(decode, "accoutType1") {

	} else {
		res.Msg = "类型错误"
		res.Name = "S_error"
		msg.Data = make(map[string]interface{})
		msg.RequestID = 1
		res.Data = msg
		raw, e := json.Marshal(res)
		if e != nil {
			fmt.Println(e)
		}
		return raw
	}

	msg.Data = make(map[string]interface{})
	msg.Put("userid", 1)
	msg.Put("token", "12345")
	msg.RequestID = 1
	res.Data = msg
	raw, e := json.Marshal(res)
	if e != nil {
		fmt.Println(e)
	}
	return raw
}
func (wb Webcenter) Register(b interface{}) interface{} {

	return nil
}

func (wb *Webcenter) RegisterController() map[string]reflect.Value {
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
