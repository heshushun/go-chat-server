package socket

import (
	"fmt"
	"go-chat/model"
	"go-chat/service/chat"
	"log"
	"strings"
)

// Hub 相当于一个事物管理中心
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// 房间号 key:client value:房间号
	roomID map[*Client]string
}

// NewHub .实例化一个hub
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		roomID:     make(map[*Client]string),
	}
}

// Run .监听客户端
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true                 // 注册client端
			h.roomID[client] = string(client.roomID) // 给client端添加房间号
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.roomID, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				// 使用“&”对message进行message切割 获取房间号
				// 向信息所属的房间内的所有client 内添加send
				// msg[0]为房间号 msg[1]为打印内容
				msg := strings.Split(string(message), "&")
				fmt.Println("房间号:" + msg[0])
				mess := strings.Split(msg[1], ":")
				fmt.Println("发言用户:", mess[0])
				fmt.Println("发言内容:", mess[1])

				storeMessage(msg[0], mess[0], mess[1])

				if string(client.roomID) == msg[0] {
					select {
					case client.send <- []byte(msg[1]):
					default:
						close(client.send)
						delete(h.clients, client)
						delete(h.roomID, client)
					}
				}
			}
		}
	}
}

// storeMessage 持久化用户的聊天消息记录
func storeMessage(roomNumber, username, content string) {
	ok, err := chat.FindRoom(roomNumber)

	if !ok {
		if err == nil {
			err := chat.CreateRoom(roomNumber)
			if err != nil {
				log.Fatalf("房间创建失败, %v", err)
			}
		} else {
			log.Fatalf("数据库查找房间失败，%v", err)
		}
	}

	var message model.Message

	message.RoomID, _ = chat.GetRoomID(roomNumber)

	message.OwnerID, _ = chat.GetUserID(username)

	message.Content = content

	err = chat.CreateMessage(&message)

	if err != nil {
		log.Fatalf("创建消息失败，%v", err)
	}
}
