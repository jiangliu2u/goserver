package game

import (
	"c-server/util"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

//实现从客户端的请求中ClientMessage中提取data(string map[string]interface{})
//name 为 NetWork.Request("webcenter",data)中的第一个参数,包括请求的模块和方法,data为本次请求所带的数据
func (cm ClientMessage) getMsg() (string, map[string]interface{}) {
	var name string
	var data map[string]interface{} //声明变量，不分配内存
	data = make(map[string]interface{})
	for one := range cm.Data {
		if one == "msg" {
			data = cm.Data[one].(map[string]interface{})
		}
		if one == "name" {
			name = cm.Data[one].(string)
		}
	}
	return name, data
}

//统一的响应客户端请求的方法
func (cm *ClientMessage) Response(msg ResponseData) {
	var res ResponseMessage
	res.Name = "S_response"
	tem := cm.Msg.(map[string]interface{})
	msg.RequestID = tem["requestID"].(float64)
	fmt.Println(cm.Msg)
	res.Data = msg
	raw, e := json.Marshal(res)
	if e != nil {
		fmt.Println(e)
		return
	}
	e = cm.Socket.socket.WriteMessage(1, raw)
	if e != nil {
		cm.Socket.WsClose()
		fmt.Println("write message error: ", e)
	}
}

//每个websocket链接的入口
func WsHandler(c *gin.Context) {
	ws, err := Upgrater.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	uid := uuid.NewV4()
	defer ws.Close()
	defer Sockets.Delete(uid)
	Sockets.Store(uid, ws)
	sock := &WSSocket{
		socket:   ws,
		isClosed: true,
	}
	go sock.HeartBeat()
	sock.On("connection", func(i interface{}) {

	})
	sock.On("C_data", func(i interface{}) {
		req := i.(ClientMessage)
		name, msg := req.getMsg() //name websocket/login
		req.Msg = msg
		head := strings.Split(name, "/")
		if len(head) < 2 {
			fmt.Println("error no such action")
			return
		}
		__moduleName := head[0] //webcenter
		temp := head[1]
		__action := "" //webcenter.ts里的action()修饰的方法名称
		//把action的首字母转化为大写,客户端是小写发送过来的
		for idx := 0; idx < len(temp); idx++ {
			if idx == 0 {
				__action = __action + strings.ToUpper(string(temp[idx]))
			} else {
				__action = __action + string(temp[idx])
			}
		}
		hehe, _ := AppControllerModule.Load(__moduleName) //比如webcenter.ts
		if hehe == nil {
			fmt.Println("error no such module")
			return
		}
		__module := hehe.(map[string]reflect.Value)
		if __module != nil {
			if util.HasFunc(__module, __action) {
				__module[__action].Call([]reflect.Value{reflect.ValueOf(req)})
			} else {
				fmt.Println("error no such action")
			}
		}
	})
	for {
		//读取ws中的数据
		_, message, err := sock.socket.ReadMessage()
		if err != nil {
			err = sock.socket.Close()
		}
		data := ClientMessage{}
		data.Socket = sock
		er := json.Unmarshal(message, &data)
		if er != nil {
			fmt.Println(err)
		}
		sock.Trigger(data.Name, data)
	}
}
