package user

import (
	"go-chat/errmsg"
	"go-chat/model"
	Service "go-chat/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Registry 用户注册
func Registry(c *gin.Context) {
	var u model.User
	_ = c.ShouldBind(&u)

	code := Service.CheckUsername(u.Username)

	if code != errmsg.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg":  errmsg.CodeMsg[code],
			"data": map[string]interface{}{
				"username": u.Username,
			},
		})

		return
	}

	var err error

	code, err = Service.AddUser(&u)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": code,
			"msg":  errmsg.CodeMsg[code],
			"data": map[string]interface{}{
				"username": u.Username,
			},
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  "注册成功",
		"data": map[string]interface{}{
			"username": u.Username,
		},
	})
}
