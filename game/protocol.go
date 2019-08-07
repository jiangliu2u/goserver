package game

import (
	"github.com/gorilla/websocket"
	"sync"
)

var eventHandlers = make(map[string]func(interface{}))
var wsEventHandlers = make(map[string]func(interface{}))

// 客户端读写消息
type wsMessage struct {
	messageType int
	data        []byte
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn    // 底层websocket
	inChan   chan ClientMessage // 读队列
	outChan  chan *wsMessage    // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}
type WSSocket struct {
	socket   *websocket.Conn
	mutex    sync.Mutex
	isClosed bool
}

//客户端传来的数据进行封装,传给每个模块
type ClientMessage struct {
	Name   string                 `json:"name"`
	Data   map[string]interface{} `json:"data"`
	Socket *wsConnection
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
//返还给客户端的data S_response 的数据
type ResponseData struct {
	Data      map[string]interface{} `json:"data"`
	RequestID float64                `json:"requestID"`
	Msg interface{} `json:"msg"`
}

//Put 往返回给客户端的数据里插入数据
func (rdi ResponseData) Put(key string, data interface{}) {
	if _, ok := rdi.Data[key]; ok {
	} else {
		rdi.Data[key] = data
	}
}

var AppControllerModule sync.Map
var Sockets sync.Map //存储
