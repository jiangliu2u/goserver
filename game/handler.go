package game

import (
	"c-server/util"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

//客户端传来的数据
type ClientMessage struct {
	Name   string                 `json:"name"`
	Data   map[string]interface{} `json:"data"`
	Socket *WSSocket
}

//返回给客户端的所有消息
type ResponseMessage struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

/**
* response{
*	name:"S_response"
*   data:any
*  }
 */
//返还给客户端的data S_response
type ResponseData struct {
	Data      map[string]interface{} `json:"data"`
	RequestID float64                `json:"requestID"`
}
type ResponseDataInterface interface {
	Put(string, interface{})
}

//Put
func (rdi ResponseData) Put(key string, data interface{}) {
	if _, ok := rdi.Data[key]; ok {

	} else {
		rdi.Data[key] = data
	}
}

//实现从ClientMessage中提取data(string map[string]interface{})
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

func WsHandler(c *gin.Context) {
	ws, err := Upgrater.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	sock := WSSocket{ws}
	sock.On("connection", func(i interface{}) {

	})
	sock.On("C_data", func(i interface{}) {
		data := i.(ClientMessage)
		name, d := data.getMsg()
		head := strings.Split(name, "/")
		actionName := head[0]
		action := head[1]
		action = "Login"
		CC := AppControllerModule[actionName]
		if util.HasFunc(CC, action) {
			back := CC[action].Call([]reflect.Value{reflect.ValueOf(d)})
			toClient := back[0].Interface().([]byte)
			e := sock.socket.WriteMessage(1, toClient)
			if e != nil {
				fmt.Println(e)
			}
		} else {
			fmt.Println("error no such action")
		}
	})
	for {
		//读取ws中的数据
		_, message, err := sock.socket.ReadMessage()
		if err != nil {
			break
		}
		data := ClientMessage{}
		data.Socket = &sock
		er := json.Unmarshal(message, &data)
		if er != nil {
			fmt.Println(err)
		}
		sock.Trigger(data.Name, data)
	}
}
