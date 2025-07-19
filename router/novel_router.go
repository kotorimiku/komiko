package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterNovelRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewNovelHandler(service)

	comicPublicGroup := publicGroup.Group("/:id/novel")
	{
		comicPublicGroup.GET("/image/*path", handler.GetImage)
		comicPublicGroup.GET("/file/*path", handler.GetFile)
	}

	comicAuthGroup := authGroup.Group("/:id/novel")
	{
		comicAuthGroup.GET("/chapters", handler.GetChapters)
		comicAuthGroup.GET("/:page", handler.GetChapter)

	}

}
