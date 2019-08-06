package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var eventHandlers = make(map[string]func(interface{}))
var wsEventHandlers = make(map[string]func(interface{}))

type eventReceiver interface {
	Trigger(string, interface{})
	On(string, func(interface{}))
}

type WSSocket struct {
	socket   *websocket.Conn
	mutex    sync.Mutex
	isClosed bool
}

func (wsConn *WSSocket) WsClose() {
	wsConn.socket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
	}
}

func (wsConn *WSSocket) HeartBeat() {
	for {
		time.Sleep(2 * time.Second)
		if err := wsConn.socket.WriteMessage(websocket.TextMessage, []byte("heartbeat from server")); err != nil {
			wsConn.WsClose()
		}
	}
}

func (ws *WSSocket) On(event string, handler func(interface{})) {
	wsEventHandlers[event] = handler
}
func (ws *WSSocket) Trigger(event string, data interface{}) {
	if len(event) == 0 {
		return
	}
	fmt.Println("trigger===========>", event)
	wsEventHandlers[event](data)
}
