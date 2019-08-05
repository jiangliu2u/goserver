package game


//客户端传来的数据
type ClientMessage struct {
	Name   string                 `json:"name"`
	Data   map[string]interface{} `json:"data"`
	Socket *WSSocket
	Msg    interface{}
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

//Put
func (rdi ResponseData) Put(key string, data interface{}) {
	if _, ok := rdi.Data[key]; ok {
	} else {
		rdi.Data[key] = data
	}
}