package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewBookHandler(service)

	bookPublicGroup := publicGroup.Group("/book")
	bookAuthGroup := authGroup.Group("/book")
	{
		RegisterComicRoutes(bookAuthGroup, bookPublicGroup, service)
		RegisterNovelRoutes(bookAuthGroup, bookPublicGroup, service)
	}

	{
		bookAuthGroup.GET("", handler.GetBooks)
		bookAuthGroup.GET("/:id", handler.GetByID)
	}

}
