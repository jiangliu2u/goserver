package game

import (
	uuid "github.com/satori/go.uuid"
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

func (wb Webcenter) Login(req ClientMessage) {
	res := ResponseData{}
	res.Data= make(map[string]interface{})
	res.Put("uid", 1)
	res.Put("token", uuid.NewV4())
	req.Error("草拟吗啊",res)
	//req.Response(res)
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
