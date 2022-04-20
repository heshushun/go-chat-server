package middlerware

import (
	"go-chat/errmsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

var MySecret = []byte("my secret")

type MyClaims struct {
	Username string `json:"username"` // 利用中间件保存一些有用的信息
	jwt.StandardClaims
}

// GenerateToken 生成token
func GenerateToken(username string) (string, int) {
	Claims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(), // 设置过期时间
			Issuer:    "peter",                               // 设置签发人
		},
	}
	reqClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	token, err := reqClaims.SignedString(MySecret)

	if err != nil {
		return "", errmsg.Error
	}

	return token, errmsg.Success
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyClaims, int) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, errmsg.Error
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, errmsg.Success
	} else {
		return nil, errmsg.InvalidToken
	}
}

// JWTAuthMiddleware jwt中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == " " {
			c.JSON(http.StatusOK, gin.H{
				"status": errmsg.CodeMsg[errmsg.AuthEmpty],
				"msg":    "请求头中的auth格式有误",
			})
			c.Abort()

			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"status": errmsg.InvalidToken,
				"msg":    errmsg.CodeMsg[errmsg.InvalidToken],
			})
			c.Abort()

			return
		}

		claims, code := ParseToken(parts[1])

		// token失效
		if code == errmsg.InvalidToken {
			c.JSON(http.StatusOK, gin.H{
				"status": errmsg.InvalidToken,
				"msg":    errmsg.CodeMsg[errmsg.InvalidToken],
			})
			c.Abort()

			return
		}
		// token过期
		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusOK, gin.H{
				"status": errmsg.TokenRunTimeError,
				"msg":    errmsg.CodeMsg[errmsg.TokenRunTimeError],
			})
			c.Abort()

			return
		}

		c.Set("username", claims.Username)

		c.Next()
	}
}
