package upload

import (
	"go-chat/errmsg"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

// PushFile 上传文件到服务器
func PushFile(c *gin.Context) {
	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": errmsg.Error,
			"msg": map[string]interface{}{
				"detail": "参数错误",
				"data":   file,
			},
		})
	}

	err = c.Request.ParseMultipartForm(100 << 20)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": errmsg.FileTooBig,
			"msg": map[string]interface{}{
				"msg":  errmsg.CodeMsg[errmsg.FileTooBig],
				"data": file,
			},
		})

		return
	}

	dist := path.Join("./sources", file.Filename)

	_ = c.SaveUploadedFile(file, dist)

	c.JSON(http.StatusOK, gin.H{
		"code": errmsg.Success,
		"msg": map[string]interface{}{
			"detail": "上传成功",
			"data":   file,
		},
	})
}

// DownLoadFile 下载文件
func DownLoadFile(c *gin.Context) {
	fileDir := c.Query("fileDir")
	fileName := c.Query("fileName")

	_, errOpenFile := os.Open(fileDir + "/" + fileName)

	if errOpenFile != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": errmsg.Error,
			"msg": map[string]interface{}{
				"detail": "文件打开错误",
			},
		})

		return
	}

	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(fileDir + "/" + fileName)
}
