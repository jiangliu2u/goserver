package game

import (
	"c-server/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsClose 关闭ws
func (wsConn *WSSocket) WsClose() {
	wsConn.socket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
	}
}

// HeartBeat 心跳
func (wsConn *WSSocket) HeartBeat() {
	for {
		time.Sleep(2 * time.Second)
		if err := wsConn.socket.WriteMessage(websocket.TextMessage, []byte("heartbeat from server")); err != nil {
			wsConn.WsClose()
		}
	}
}

//Init 初始化 websocket模块 注册对应的游戏路由
func Init() {
	//AppControllerModule = make(map[string]map[string]reflect.Value)
	wc := &Webcenter{"webcenter"}
	wb := &Wcenter{"wcenter"}
	AppControllerModule.Store(wc.Name, RegisterController(wc))
	AppControllerModule.Store(wb.Name, RegisterController(wb))
}

//事件监听
func (wsConn *WsConnection) Listen(event string, handler func(message ClientMessage)) {
	wsEventHandlers[event] = handler
}

//事件分发
func (wsConn *WsConnection) Dispatch(event string, data ClientMessage) {
	if len(event) == 0 {
		return
	}
	wsEventHandlers[event](data)
}

//Handler ws处理
func Handler(c *gin.Context) {
	wsSocket, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsConn := &WsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan ClientMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	// uid := util.RandStringRunes(11)
	WSClients["uid"] = wsConn
	// 处理器
	// go wsConn.procLoop()
	// 读协程
	wsConn.Listen("C_data", func(req ClientMessage) {
		//req := i.(ClientMessage)
		name, msg := req.getMsg() //name websocket/login
		req.Msg = msg
		head := strings.Split(name, "/")
		if len(head) < 2 {
			fmt.Println("error no such action")
			return
		}
		moduleName := head[0] //webcenter
		temp := head[1]
		action := "" //webcenter.ts里的action()修饰的方法名称
		//把action的首字母转化为大写,客户端是小写发送过来的
		for idx := 0; idx < len(temp); idx++ {
			if idx == 0 {
				action = action + strings.ToUpper(string(temp[idx]))
			} else {
				action = action + string(temp[idx])
			}
		}
		hehe, _ := AppControllerModule.Load(moduleName) //比如webcenter.ts
		if hehe == nil {
			req.error("模块" + moduleName + "不存在...")
			fmt.Println("模块" + moduleName + "不存在...")
			return
		}
		module := hehe.(map[string]reflect.Value)
		fmt.Println("模块" + moduleName + "的方法" + action)
		if module != nil {
			if util.HasFunc(module, action) {
				module[action].Call([]reflect.Value{reflect.ValueOf(req)})
			} else {
				req.error("模块" + moduleName + "的方法" + action + "不存在...")
				fmt.Println("模块" + moduleName + "的方法" + action + "不存在...")
			}
		}
	})
	//读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
	//处理收到的消息 协程
	go wsConn.handRecievedMsg()
}

func (wsConn *WsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

func (wsConn *WsConnection) wsRead() (data interface{}, err error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

func (wsConn *WsConnection) wsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

func (wsConn *WsConnection) handRecievedMsg() {
	for {
		select {
		case msg := <-wsConn.inChan:
			wsConn.Dispatch(msg.Name, msg)
		case <-wsConn.closeChan:
			return
		}
	}
}

func (wsConn *WsConnection) wsReadLoop() {
	for {
		_, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			wsConn.wsClose()
			break
		}
		//data ,err := wsConn.wsRead()
		// if err != nil {
		// 	wsConn.wsClose()
		// 	break
		// }
		cmg := ClientMessage{}
		cmg.Socket = wsConn
		er := json.Unmarshal(data, &cmg)
		if er != nil {
			wsConn.wsClose()
			break
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- cmg:
		case <-wsConn.closeChan:
			break
		}
	}

}

func (wsConn *WsConnection) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				wsConn.wsClose()
				break
			}
		case <-wsConn.closeChan:
			break
		}
	}

}
