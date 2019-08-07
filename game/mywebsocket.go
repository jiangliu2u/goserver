package game

import (
	"time"
	"github.com/gorilla/websocket"
)

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

