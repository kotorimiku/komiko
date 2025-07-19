package router

import (
	"komiko/handler"

	"github.com/gin-gonic/gin"
)

func RegisterCoverRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup) {

	coverGroup := publicGroup.Group("/cover")
	{
		coverGroup.GET("/:file", handler.GetCover)
	}

}
