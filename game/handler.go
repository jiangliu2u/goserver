package game

import (
	"encoding/json"
	"fmt"
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

//统一的响应客户端请求的方法 每个模块里的返回
func (cm *ClientMessage) Response(msg ResponseData) {
	var res ResponseMessage
	res.Name = "S_response"
	tem := cm.Msg.(map[string]interface{})
	msg.RequestID = tem["requestID"].(float64)
	res.Data = msg
	raw, e := json.Marshal(res)
	if e != nil {
		fmt.Println(e)
		return
	}
	e = cm.Socket.wsWrite(1, raw)
	if e != nil {
		fmt.Println(e)
		return
	}
}

//统一的响应客户端请求的方法 每个模块里的返回
func (cm *ClientMessage) Error(info string, msg ResponseData) {
	var res ResponseMessage
	res.Name = "S_error"
	tem := cm.Msg.(map[string]interface{})
	msg.RequestID = tem["requestID"].(float64)
	msg.Msg=info
	res.Data = msg
	raw, e := json.Marshal(res)
	if e != nil {
		fmt.Println(e)
		return
	}
	e = cm.Socket.wsWrite(1, raw)
	if e != nil {
		fmt.Println(e)
		return
	}
}
