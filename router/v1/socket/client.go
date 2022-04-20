package socket

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// 用户名称
	username []byte
	// 房间号
	roomID []byte
}

// 从时间处理中心读取消息
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { _ = c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		message = bytes.TrimSpace(bytes.ReplaceAll(message, newline, space))

		message = []byte(string(c.roomID) + "&" + string(c.username) + ":" + string(message))
		fmt.Println(string(message))
		c.hub.broadcast <- message
	}
}

// 从事物处理中心读取消息
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		_ = c.conn.Close()

	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// 使用“&”分割获取房间号
			// 聊天内容不得包含&字符
			// msg[0]为房间号 msg[1]为打印内容
			// msg := strings.Split(string(message), "&")
			// if msg[0] == string(c.hub.roomID[c]) {
			// 	w.Write([]byte(msg[1]))
			// }
			_, _ = w.Write(message)
			// Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	if msg[0] == string(c.hub.roomID[c]) {
			// 		w.Write(newline)
			// 		w.Write(<-c.send)
			// 	}
			// }
			if err := w.Close(); err != nil {
				log.Printf("error: %v", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ChatRequest 聊天室请求
type ChatRequest struct {
	RoomID   string `json:"room_id" form:"room_id"`
	UserID   int    `json:"user_id" form:"user_id"`
	UserName string `json:"user_name" form:"user_name"`
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, c *gin.Context) {
	// 获取前端数据
	var req ChatRequest

	_ = c.ShouldBind(&req)

	userName := req.UserName
	roomID := req.RoomID
	// 获取redis连接(暂未使用)
	// pool := c.MustGet("test").(*redis.Pool)
	// redisConn := pool.Get()
	// defer redisConn.Close()
	// 将网络请求变为websocket
	var upGrader = websocket.Upgrader{
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println(err)
		return
	}

	// fmt.Println("username" + userName)

	// fmt.Println("roomID" + roomID)

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: []byte(userName),
		roomID:   []byte(roomID),
	}

	client.hub.register <- client

	go client.writePump()

	go client.readPump()
}
