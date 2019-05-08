package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"goApiServer/ws"
	"net/http"
	"time"
)

type API struct {
	beego.Controller
}

type SendMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type WebSocketController struct {
	beego.Controller
}

var (
	// Configure the upgrader
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (c *WebSocketController) SayHi() {
	beego.Notice("收到请求。。。")

	var (
		wsConn *websocket.Conn
		err    error
		// data []byte
		conn *ws.Connection
		data []byte
	)

	// 完成http应答
	if wsConn, err = upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil); err != nil {
		return // 获取连接失败直接返回
	}

	if conn, err = ws.InitConnection(wsConn); err != nil {
		goto ERR
	}

	go func() {
		var (
			err error
		)
		for {
			// 每隔一秒发送一次心跳
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}

	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		beego.Debug("收到一条消息:", string(data))
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	// 关闭连接
	conn.Close()
}
