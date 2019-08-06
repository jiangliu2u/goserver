package game

import (
	"c-server/util"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 客户端读写消息
type wsMessage struct {
	messageType int
	data        []byte
}

// 客户端连接
type wsConnection struct {
	wsSocket *websocket.Conn     // 底层websocket
	inChan   chan *ClientMessage // 读队列
	outChan  chan *wsMessage     // 写队列

	mutex     sync.Mutex // 避免重复关闭管道
	isClosed  bool
	closeChan chan byte // 关闭通知
}

func (ws *wsConnection) Listen(event string, handler func(interface{})) {
	wsEventHandlers[event] = handler
}
func (ws *wsConnection) Dispatch(event string, data interface{}) {
	if len(event) == 0 {
		return
	}
	fmt.Println("trigger===========>", event)
	wsEventHandlers[event](data)
}

func Handler(c *gin.Context) {
	wsSocket, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *ClientMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}
	// 处理器
	// go wsConn.procLoop()
	// 读协程
	wsConn.Listen("C_data", func(i interface{}) {
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
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}

func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}

func (wsConn *wsConnection) wsRead() (*ClientMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}

func (wsConn *wsConnection) wsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}

func (wsConn *wsConnection) wsReadLoop() {
	for {
		// 读一个message
		/*
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
		*/

		_, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			goto error
		}
		cmg := &ClientMessage{}
		cmg.Socket = wsConn
		er := json.Unmarshal(data, &cmg)
		if er != nil {
			fmt.Println(err)
		}
		// 放入请求队列
		select {
		case wsConn.inChan <- cmg:
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}

func (wsConn *wsConnection) wsWriteLoop() {
	for {
		select {
		// 取一个应答
		case msg := <-wsConn.outChan:
			// 写给websocket
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				goto error
			}
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}
