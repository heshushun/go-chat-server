package user

import (
	"go-chat/errmsg"
	"go-chat/middlerware"
	"go-chat/model"
	Service "go-chat/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthHandler 登录验证
func AuthHandler(c *gin.Context) {
	var user model.User

	err := c.ShouldBind(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": errmsg.Error,
			"msg": map[string]interface{}{
				"detail": "无效的参数",
				"data":   "",
			},
		})

		return
	}

	code, err := Service.CheckLogin(&user)

	if code != errmsg.Success {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": code,
			"msg": map[string]interface{}{
				"status": errmsg.CodeMsg[code],
			},
		})
	}

	tokenString, code := middlerware.GenerateToken(user.Username)

	c.JSON(http.StatusOK, gin.H{
		"status": code,
		"msg":    "登录成功",
		"data": map[string]interface{}{
			"token":    tokenString,
			"username": user.Username,
		},
	})
}
