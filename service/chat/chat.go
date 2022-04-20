package chat

import (
	"go-chat/database"
	"go-chat/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// GetUserID 通过用户名获取用户的ID
func GetUserID(username string) (int, error) {
	var user model.User
	if err := database.Db.
		Where("username = ?", username).
		First(&user).
		Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

// GetRoomID 通过房间号获取房间的ID
func GetRoomID(number string) (int, error) {
	var room model.Room
	if err := database.Db.
		Where("room_number = ?", number).
		First(&room).
		Error; err != nil {
		return 0, err
	}

	return room.ID, nil
}

// CreateMessage 创建聊天信息
func CreateMessage(message *model.Message) error {
	return database.Db.
		Create(message).
		Error
}

// CreateRoom 创建一个新的房间
func CreateRoom(number string) error {
	var room model.Room
	room.RoomNumber = number
	return database.Db.
		Create(&room).Error
}

// FindRoom 通过房间号查找房间是否存在，存在返回ture，不存在返回false
func FindRoom(number string) (bool, error) {
	var room model.Room

	if err := database.Db.Where("room_number = ?", number).First(&room).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
