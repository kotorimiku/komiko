package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterComicRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewComicHandler(service)

	comicGroup := publicGroup.Group("/:id/comic")
	{
		comicGroup.GET("/:page", handler.GetImage)
	}

}
