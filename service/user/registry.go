package user

import (
	"encoding/base64"
	"go-chat/database"
	"go-chat/errmsg"
	"go-chat/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/scrypt"
	"log"
)

// CheckUsername 检查用户名是否存在
func CheckUsername(username string) (code int) {
	var user model.User
	if err := database.Db.
		Where("username = ?", username).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errmsg.Success
		}

		return errmsg.ErrMysqlServer
	}

	return errmsg.ErrUsernameExist
}

// AddUser 添加用户
func AddUser(user *model.User) (code int, err error) {
	user.Password = ScryptPassword(user.Password)
	if err := database.Db.
		Create(user).
		Error; err != nil {
		return errmsg.ErrMysqlServer, err
	}

	return errmsg.Success, nil
}

// ScryptPassword 对密码进行加密
func ScryptPassword(password string) string {
	const keyLen = 10
	salt := make([]byte, 8)
	salt = []byte{23, 21, 34, 7, 29, 12, 34, 0}

	HashPassword, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, keyLen)

	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(HashPassword)
}
