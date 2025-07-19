package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": data,
	})
}

func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": 500,
		"msg":  msg,
	})
}
