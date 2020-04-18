package game

import (
	"sync"

	"github.com/gorilla/websocket"
)

var eventHandlers = make(map[string]func(interface{}))
var wsEventHandlers = make(map[string]func(message ClientMessage))

//WSClients ws的map
var WSClients map[string]*WsConnection

// 客户端读写消息
type wsMessage struct {
	messageType int
	data        []byte
}

// 客户端连接
type WsConnection struct {
	wsSocket *websocket.Conn    // 底层websocket
	inChan   chan ClientMessage // 读队列
	outChan  chan *wsMessage    // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

// WSSocket ws封装
type WSSocket struct {
	socket   *websocket.Conn
	mutex    sync.Mutex
	isClosed bool
}

//ClientMessage 客户端传来的数据进行封装,传给每个模块
type ClientMessage struct {
	Name   string                 `json:"name"`
	Data   map[string]interface{} `json:"data"`
	Socket *WsConnection
	Msg    interface{}
}

//ResponseMessage 返回给客户端的所有消息
type ResponseMessage struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

//ResponseData 返还给客户端的data S_response 的数据
type ResponseData struct {
	Data      map[string]interface{} `json:"data"`
	RequestID float64                `json:"requestID"`
	Msg       interface{}            `json:"msg"`
}

//Put 往返回给客户端的数据里插入数据
func (rdi ResponseData) Put(key string, data interface{}) {
	if _, ok := rdi.Data[key]; ok {
	} else {
		rdi.Data[key] = data
	}
}

//AppControllerModule 存储模块
var AppControllerModule sync.Map

//Sockets 存储ws
var Sockets sync.Map
