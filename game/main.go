package game

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var Upgrater = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var AppControllerModule map[string]map[string]reflect.Value

//func registerController(contro *interface{}) map[string]reflect.Value {
//	cont := make(map[string]reflect.Value)
//	v := reflect.ValueOf(&contro).Elem()
//	t := reflect.TypeOf(&contro).Elem()
//	for i := 0; i < v.NumMethod(); i++ {
//		action := v.Method(i)
//		actionName := t.Method(i).Name
//		cont[actionName] = action
//	}
//	return cont
//}

func Init() {
	AppControllerModule = make(map[string]map[string]reflect.Value)
	wc := &Webcenter{"webcenter"}
	AppControllerModule[wc.Name] = wc.RegisterController()
}

//所有websocket的入口
func Entry(c *gin.Context) {
	go WsHandler(c) //每个ws开一个goroutine
}
