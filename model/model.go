package model

import "time"

// User 用户表
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Room 房间表
type Room struct {
	ID         int    `json:"id"`
	RoomNumber string `json:"room_number"`
}

// Message 消息表
type Message struct {
	ID        int        `json:"id"`
	OwnerID   int        `json:"owner_id"`
	RoomID    int        `json:"room_id"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	Content   string     `json:"content"`
}
