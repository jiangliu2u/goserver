package game

import (
	"fmt"
	"github.com/gorilla/websocket"
)

var eventHandlers = make(map[string]func(interface{}))
var wsEventHandlers = make(map[string]func(interface{}))

type eventReceiver interface {
	Trigger(string, interface{})
	On(string, func(interface{}))
}

type WSSocket struct {
	socket *websocket.Conn
}

func (ws *WSSocket) On(event string, handler func(interface{})) {
	wsEventHandlers[event] = handler
}
func (ws *WSSocket) Trigger(event string, data interface{}) {
	fmt.Println("trigger===========>",event)
	wsEventHandlers[event](data)
}
