package router

import (
	"komiko/handler"
	"komiko/service"

	"github.com/gin-gonic/gin"
)

func RegisterProgressRoutes(authGroup *gin.RouterGroup, publicGroup *gin.RouterGroup, service *service.Service) {

	handler := handler.NewProgressHandler(service)

	progressGroup := authGroup.Group("/progress")
	{
		progressGroup.GET("", handler.GetBookProgresses)
		progressGroup.GET("/series", handler.GetSeriesProgresses)
		progressGroup.GET("/:id", handler.GetByID)
		progressGroup.POST("", handler.Create)
	}

}
